/**
 * Teleport
 * Copyright (C) 2025 Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

import { ResourceKind } from 'teleport/services/resources';

import {
  ResourceAccess,
  roleToRoleEditorModel,
  RuleModel,
} from './standardmodel';
import {
  validateAccessRule,
  validateResourceAccess,
  validateRoleEditorModel,
} from './validation';
import { withDefaults } from './withDefaults';

const minimalRoleModel = () =>
  roleToRoleEditorModel(
    withDefaults({ metadata: { name: 'role-name' }, version: 'v7' })
  );

const validity = (arr: { valid: boolean }[]) => arr.map(it => it.valid);

describe('validateRoleEditorModel', () => {
  test('valid minimal model', () => {
    const result = validateRoleEditorModel(
      minimalRoleModel(),
      undefined,
      undefined
    );
    expect(result.metadata.valid).toBe(true);
    expect(result.resources).toEqual([]);
    expect(result.rules).toEqual([]);
  });

  test('valid complex model', () => {
    const model = minimalRoleModel();
    model.metadata.labels = [{ name: 'foo', value: 'bar' }];
    model.resources = [
      {
        kind: 'kube_cluster',
        labels: [{ name: 'foo', value: 'bar' }],
        groups: [],
        resources: [
          {
            id: 'dummy-id',
            kind: { label: 'pod', value: 'pod' },
            name: 'res-name',
            namespace: 'dummy-namespace',
            verbs: [],
          },
        ],
      },
      {
        kind: 'node',
        labels: [{ name: 'foo', value: 'bar' }],
        logins: [{ label: 'root', value: 'root' }],
      },
      {
        kind: 'app',
        labels: [{ name: 'foo', value: 'bar' }],
        awsRoleARNs: ['some-arn'],
        azureIdentities: ['some-azure-id'],
        gcpServiceAccounts: ['some-gcp-acct'],
      },
      {
        kind: 'db',
        labels: [{ name: 'foo', value: 'bar' }],
        roles: [{ label: 'some-role', value: 'some-role' }],
        names: [],
        users: [],
      },
      {
        kind: 'windows_desktop',
        labels: [{ name: 'foo', value: 'bar' }],
        logins: [],
      },
    ];
    model.rules = [
      {
        id: 'dummy-id',
        resources: [{ label: ResourceKind.Node, value: ResourceKind.Node }],
        verbs: [{ label: '*', value: '*' }],
      },
    ];
    const result = validateRoleEditorModel(model, undefined, undefined);
    expect(result.metadata.valid).toBe(true);
    expect(validity(result.resources)).toEqual([true, true, true, true, true]);
    expect(validity(result.rules)).toEqual([true]);
  });

  test('invalid metadata', () => {
    const model = minimalRoleModel();
    model.metadata.name = '';
    const result = validateRoleEditorModel(model, undefined, undefined);
    expect(result.metadata.valid).toBe(false);
  });

  test('invalid resource', () => {
    const model = minimalRoleModel();
    model.resources = [
      {
        kind: 'node',
        labels: [{ name: 'foo', value: '' }],
        logins: [],
      },
    ];
    const result = validateRoleEditorModel(model, undefined, undefined);
    expect(validity(result.resources)).toEqual([false]);
  });

  test('invalid access rule', () => {
    const model = minimalRoleModel();
    model.rules = [
      {
        id: 'dummy-id',
        resources: [],
        verbs: [{ label: '*', value: '*' }],
      },
    ];
    const result = validateRoleEditorModel(model, undefined, undefined);
    expect(validity(result.rules)).toEqual([false]);
  });

  it('reuses previously computed section results', () => {
    const model = minimalRoleModel();
    const result1 = validateRoleEditorModel(model, undefined, undefined);
    const result2 = validateRoleEditorModel(model, model, result1);
    expect(result2.metadata).toBe(result1.metadata);
    expect(result2.resources).toBe(result1.resources);
    expect(result2.rules).toBe(result1.rules);
  });
});

describe('validateResourceAccess', () => {
  it('reuses previously computed results', () => {
    const resource: ResourceAccess = { kind: 'node', labels: [], logins: [] };
    const result1 = validateResourceAccess(resource, undefined, undefined);
    const result2 = validateResourceAccess(resource, resource, result1);
    expect(result2).toBe(result1);
  });
});

describe('validateAccessRule', () => {
  it('reuses previously computed results', () => {
    const rule: RuleModel = { id: 'some-id', resources: [], verbs: [] };
    const result1 = validateAccessRule(rule, undefined, undefined);
    const result2 = validateAccessRule(rule, rule, result1);
    expect(result2).toBe(result1);
  });
});
