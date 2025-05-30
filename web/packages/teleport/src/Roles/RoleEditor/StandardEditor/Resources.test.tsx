/**
 * Teleport
 * Copyright (C) 2024 Gravitational, Inc.
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

import { act, within } from '@testing-library/react';
import selectEvent from 'react-select-event';

import { render, screen, userEvent } from 'design/utils/testing';
import { Validator } from 'shared/components/Validation';

import { RoleVersion } from 'teleport/services/resources';

import {
  AppAccessSection,
  DatabaseAccessSection,
  GitHubOrganizationAccessSection,
  KubernetesAccessSection,
  ServerAccessSection,
  WindowsDesktopAccessSection,
} from './Resources';
import {
  AppAccess,
  DatabaseAccess,
  defaultRoleVersion,
  GitHubOrganizationAccess,
  KubernetesAccess,
  newResourceAccess,
  ServerAccess,
  WindowsDesktopAccess,
} from './standardmodel';
import { StatefulSection } from './StatefulSection';
import {
  ResourceAccessValidationResult,
  validateResourceAccess,
} from './validation';

describe('ServerAccessSection', () => {
  const setup = () => {
    const onChange = jest.fn();
    let validator: Validator;
    render(
      <StatefulSection<ServerAccess, ResourceAccessValidationResult>
        component={ServerAccessSection}
        defaultValue={newResourceAccess('node', defaultRoleVersion)}
        onChange={onChange}
        validatorRef={v => {
          validator = v;
        }}
        validate={validateResourceAccess}
      />
    );
    return { user: userEvent.setup(), onChange, validator };
  };

  test('editing', async () => {
    const { user, onChange } = setup();
    await user.type(screen.getByPlaceholderText('label key'), 'some-key');
    await user.type(screen.getByPlaceholderText('label value'), 'some-value');
    await selectEvent.create(screen.getByLabelText('Logins'), 'root', {
      createOptionText: 'Login: root',
    });
    await selectEvent.create(screen.getByLabelText('Logins'), 'some-user', {
      createOptionText: 'Login: some-user',
    });

    expect(onChange).toHaveBeenLastCalledWith({
      kind: 'node',
      labels: [{ name: 'some-key', value: 'some-value' }],
      logins: [
        expect.objectContaining({
          label: '{{internal.logins}}',
          value: '{{internal.logins}}',
        }),
        expect.objectContaining({ label: 'root', value: 'root' }),
        expect.objectContaining({ label: 'some-user', value: 'some-user' }),
      ],
      hideValidationErrors: true,
    } as ServerAccess);
  });

  test('validation', async () => {
    const { user, validator } = setup();
    await user.type(screen.getByPlaceholderText('label value'), 'some-value');
    await selectEvent.create(screen.getByLabelText('Logins'), '*', {
      createOptionText: 'Login: *',
    });
    act(() => validator.validate());
    expect(
      screen.getByPlaceholderText('label key')
    ).toHaveAccessibleDescription('required');
    expect(
      screen.getByText('Wildcard is not allowed in logins')
    ).toBeInTheDocument();
  });
});

describe('KubernetesAccessSection', () => {
  const setup = (roleVersion: RoleVersion = defaultRoleVersion) => {
    const onChange = jest.fn();
    let validator: Validator;
    render(
      <StatefulSection<KubernetesAccess, ResourceAccessValidationResult>
        component={KubernetesAccessSection}
        defaultValue={{
          ...newResourceAccess('kube_cluster', defaultRoleVersion),
          roleVersion,
        }}
        onChange={onChange}
        validatorRef={v => {
          validator = v;
        }}
        validate={validateResourceAccess}
      />
    );
    return { user: userEvent.setup(), onChange, validator };
  };

  test('editing', async () => {
    const { user, onChange } = setup();

    await selectEvent.create(screen.getByLabelText('Groups'), 'group1', {
      createOptionText: 'Group: group1',
    });
    await selectEvent.create(screen.getByLabelText('Groups'), 'group2', {
      createOptionText: 'Group: group2',
    });

    await user.type(screen.getByPlaceholderText('label key'), 'some-key');
    await user.type(screen.getByPlaceholderText('label value'), 'some-value');

    await selectEvent.create(screen.getByLabelText('Users'), 'joe', {
      createOptionText: 'User: joe',
    });
    await selectEvent.create(screen.getByLabelText('Users'), 'mary', {
      createOptionText: 'User: mary',
    });

    await user.click(
      screen.getByRole('button', { name: 'Add a Kubernetes Resource' })
    );
    expect(
      reactSelectValueContainer(screen.getByLabelText('Kind (plural)'))
    ).toHaveTextContent('Any kind');
    expect(screen.getByLabelText('API Group *')).toHaveValue('*');
    expect(screen.getByLabelText('Name *')).toHaveValue('*');
    expect(screen.getByLabelText('Namespace *')).toHaveValue('*');
    await selectEvent.select(screen.getByLabelText('Kind (plural)'), 'jobs');
    await user.clear(screen.getByLabelText('API Group *'));
    await user.type(screen.getByLabelText('API Group *'), 'api-group-name');
    await user.clear(screen.getByLabelText('Name *'));
    await user.type(screen.getByLabelText('Name *'), 'job-name');
    await user.clear(screen.getByLabelText('Namespace *'));
    await user.type(screen.getByLabelText('Namespace *'), 'job-namespace');
    await selectEvent.select(screen.getByLabelText('Verbs'), [
      'create',
      'delete',
    ]);

    expect(onChange).toHaveBeenLastCalledWith({
      kind: 'kube_cluster',
      groups: [
        expect.objectContaining({ value: '{{internal.kubernetes_groups}}' }),
        expect.objectContaining({ value: 'group1' }),
        expect.objectContaining({ value: 'group2' }),
      ],
      labels: [{ name: 'some-key', value: 'some-value' }],
      resources: [
        {
          id: expect.any(String),
          kind: expect.objectContaining({ value: 'jobs' }),
          name: 'job-name',
          namespace: 'job-namespace',
          verbs: [
            expect.objectContaining({ value: 'create' }),
            expect.objectContaining({ value: 'delete' }),
          ],
          apiGroup: 'api-group-name',
          roleVersion: 'v8',
        },
      ],
      users: [
        expect.objectContaining({ value: 'joe' }),
        expect.objectContaining({ value: 'mary' }),
      ],
      roleVersion: 'v8',
      hideValidationErrors: true,
    } as KubernetesAccess);
  });

  test('adding and removing resources', async () => {
    const { user, onChange } = setup();

    await user.click(
      screen.getByRole('button', { name: 'Add a Kubernetes Resource' })
    );
    await user.clear(screen.getByLabelText('Name *'));
    await user.type(screen.getByLabelText('Name *'), 'res1');
    await user.click(
      screen.getByRole('button', { name: 'Add Another Kubernetes Resource' })
    );
    await user.clear(screen.getAllByLabelText('Name *')[1]);
    await user.type(screen.getAllByLabelText('Name *')[1], 'res2');
    await user.click(
      screen.getByRole('button', { name: 'Add Another Kubernetes Resource' })
    );
    await user.clear(screen.getAllByLabelText('Name *')[2]);
    await user.type(screen.getAllByLabelText('Name *')[2], 'res3');
    expect(onChange).toHaveBeenLastCalledWith(
      expect.objectContaining({
        resources: [
          expect.objectContaining({ name: 'res1' }),
          expect.objectContaining({ name: 'res2' }),
          expect.objectContaining({ name: 'res3' }),
        ],
      })
    );

    await user.click(
      screen.getAllByRole('button', { name: 'Remove Kubernetes resource' })[1]
    );
    expect(onChange).toHaveBeenLastCalledWith(
      expect.objectContaining({
        resources: [
          expect.objectContaining({ name: 'res1' }),
          expect.objectContaining({ name: 'res3' }),
        ],
      })
    );
    await user.click(
      screen.getAllByRole('button', { name: 'Remove Kubernetes resource' })[0]
    );
    expect(onChange).toHaveBeenLastCalledWith(
      expect.objectContaining({
        resources: [expect.objectContaining({ name: 'res3' })],
      })
    );
    await user.click(
      screen.getAllByRole('button', { name: 'Remove Kubernetes resource' })[0]
    );
    expect(onChange).toHaveBeenLastCalledWith(
      expect.objectContaining({ resources: [] })
    );
  });

  test('validation', async () => {
    const { user, validator } = setup(RoleVersion.V6);
    await user.type(screen.getByPlaceholderText('label value'), 'some-value');
    await user.click(
      screen.getByRole('button', { name: 'Add a Kubernetes Resource' })
    );

    screen.getByLabelText('Kind').setAttribute('value', 'Service');

    screen.getByLabelText('Kind');
    await user.clear(screen.getByLabelText('Name *'));
    await user.clear(screen.getByLabelText('Namespace *'));
    await selectEvent.select(screen.getByLabelText('Verbs'), [
      'All verbs',
      'create',
    ]);
    act(() => validator.validate());
    expect(
      screen.getByText('Only pods are allowed for role version v6')
    ).toBeVisible();
    expect(
      screen.getByPlaceholderText('label key')
    ).toHaveAccessibleDescription('required');
    expect(screen.getByLabelText('Name *')).toHaveAccessibleDescription(
      'Resource name is required, use "*" for any resource'
    );
    expect(screen.getByLabelText('Namespace *')).toHaveAccessibleDescription(
      'Namespace is required for resources of this kind'
    );
    expect(
      screen.getByText('Mixing "All verbs" with other options is not allowed')
    ).toBeVisible();
  });
});

describe('AppAccessSection', () => {
  const setup = (model: Partial<AppAccess> = {}) => {
    const onChange = jest.fn();
    let validator: Validator;
    render(
      <StatefulSection<AppAccess, ResourceAccessValidationResult>
        component={AppAccessSection}
        defaultValue={{
          ...newResourceAccess('app', defaultRoleVersion),
          ...model,
        }}
        onChange={onChange}
        validatorRef={v => {
          validator = v;
        }}
        validate={validateResourceAccess}
      />
    );
    return { user: userEvent.setup(), onChange, validator };
  };

  const awsRoleArns = () =>
    screen.getByRole('group', { name: 'AWS Role ARNs' });
  const awsRoleArnTextBoxes = () =>
    within(awsRoleArns()).getAllByRole('textbox');
  const azureIdentities = () =>
    screen.getByRole('group', { name: 'Azure Identities' });
  const azureIdentityTextBoxes = () =>
    within(azureIdentities()).getAllByRole('textbox');
  const gcpServiceAccounts = () =>
    screen.getByRole('group', { name: 'GCP Service Accounts' });
  const gcpServiceAccountTextBoxes = () =>
    within(gcpServiceAccounts()).getAllByRole('textbox');
  const mcpTools = () => screen.getByRole('group', { name: 'MCP Tools' });
  const mcpToolsTextBoxes = () => within(mcpTools()).getAllByRole('textbox');

  test('editing', async () => {
    const { user, onChange } = setup();
    await user.type(screen.getByPlaceholderText('label key'), 'env');
    await user.type(screen.getByPlaceholderText('label value'), 'prod');

    // Instead of typing these ungodly long values, we paste them — otherwise,
    // this test may time out. And that's what our users would typically do,
    // anyway.
    await user.click(
      within(awsRoleArns()).getByRole('button', { name: 'Add More' })
    );
    await user.click(awsRoleArnTextBoxes()[1]);
    await user.paste('arn:aws:iam::123456789012:role/admin');
    await user.click(
      within(azureIdentities()).getByRole('button', { name: 'Add More' })
    );
    await user.click(azureIdentityTextBoxes()[1]);
    await user.paste(
      '/subscriptions/1020304050607-cafe-8090-a0b0c0d0e0f0/resourceGroups/example-resource-group/providers/Microsoft.ManagedIdentity/userAssignedIdentities/admin'
    );
    await user.click(
      within(gcpServiceAccounts()).getByRole('button', { name: 'Add More' })
    );
    await user.click(gcpServiceAccountTextBoxes()[1]);
    await user.paste('admin@some-project.iam.gserviceaccount.com');
    await user.click(
      within(mcpTools()).getByRole('button', { name: 'Add More' })
    );
    await user.click(mcpToolsTextBoxes()[1]);
    await user.paste('allow_tools_with_prefix_*');

    expect(onChange).toHaveBeenLastCalledWith({
      kind: 'app',
      labels: [{ name: 'env', value: 'prod' }],
      awsRoleARNs: [
        '{{internal.aws_role_arns}}',
        'arn:aws:iam::123456789012:role/admin',
      ],
      azureIdentities: [
        '{{internal.azure_identities}}',
        '/subscriptions/1020304050607-cafe-8090-a0b0c0d0e0f0/resourceGroups/example-resource-group/providers/Microsoft.ManagedIdentity/userAssignedIdentities/admin',
      ],
      gcpServiceAccounts: [
        '{{internal.gcp_service_accounts}}',
        'admin@some-project.iam.gserviceaccount.com',
      ],
      mcpTools: ['{{internal.mcp_tools}}', 'allow_tools_with_prefix_*'],
      hideValidationErrors: true,
    } as AppAccess);
  });

  test('validation', async () => {
    const { user, validator } = setup();
    await user.type(screen.getByPlaceholderText('label value'), 'some-value');
    await user.click(
      within(awsRoleArns()).getByRole('button', { name: 'Add More' })
    );
    await user.type(awsRoleArnTextBoxes()[1], '*');
    await user.click(
      within(azureIdentities()).getByRole('button', { name: 'Add More' })
    );
    await user.type(azureIdentityTextBoxes()[1], '*');
    await user.click(
      within(gcpServiceAccounts()).getByRole('button', { name: 'Add More' })
    );
    await user.type(gcpServiceAccountTextBoxes()[1], '*');
    act(() => validator.validate());
    expect(
      screen.getByPlaceholderText('label key')
    ).toHaveAccessibleDescription('required');
    expect(awsRoleArnTextBoxes()[1]).toHaveAccessibleDescription(
      'Wildcard is not allowed in AWS role ARNs'
    );
    expect(azureIdentityTextBoxes()[1]).toHaveAccessibleDescription(
      'Wildcard is not allowed in Azure identities'
    );
    expect(gcpServiceAccountTextBoxes()[1]).toHaveAccessibleDescription(
      'Wildcard is not allowed in GCP service accounts'
    );
  });
});

describe('DatabaseAccessSection', () => {
  const setup = () => {
    const onChange = jest.fn();
    let validator: Validator;
    render(
      <StatefulSection<DatabaseAccess, ResourceAccessValidationResult>
        component={DatabaseAccessSection}
        defaultValue={newResourceAccess('db', defaultRoleVersion)}
        onChange={onChange}
        validatorRef={v => {
          validator = v;
        }}
        validate={validateResourceAccess}
      />
    );
    return { user: userEvent.setup(), onChange, validator };
  };

  test('editing', async () => {
    const { user, onChange } = setup();

    const labels = within(screen.getByRole('group', { name: 'Labels' }));
    await user.type(labels.getByPlaceholderText('label key'), 'env');
    await user.type(labels.getByPlaceholderText('label value'), 'prod');

    await selectEvent.create(screen.getByLabelText('Database Names'), 'stuff', {
      createOptionText: 'Database Name: stuff',
    });
    await selectEvent.create(screen.getByLabelText('Database Users'), 'mary', {
      createOptionText: 'Database User: mary',
    });
    await selectEvent.create(screen.getByLabelText('Database Roles'), 'admin', {
      createOptionText: 'Database Role: admin',
    });

    const dbServiceLabels = within(
      screen.getByRole('group', { name: 'Database Service Labels' })
    );
    await user.type(dbServiceLabels.getByPlaceholderText('label key'), 'foo');
    await user.type(dbServiceLabels.getByPlaceholderText('label value'), 'bar');

    expect(onChange).toHaveBeenLastCalledWith({
      kind: 'db',
      labels: [{ name: 'env', value: 'prod' }],
      names: [
        expect.objectContaining({ value: '{{internal.db_names}}' }),
        expect.objectContaining({ label: 'stuff', value: 'stuff' }),
      ],
      roles: [
        expect.objectContaining({ value: '{{internal.db_roles}}' }),
        expect.objectContaining({ label: 'admin', value: 'admin' }),
      ],
      users: [
        expect.objectContaining({ value: '{{internal.db_users}}' }),
        expect.objectContaining({ label: 'mary', value: 'mary' }),
      ],
      dbServiceLabels: [{ name: 'foo', value: 'bar' }],
      hideValidationErrors: true,
    } as DatabaseAccess);
  });

  test('validation', async () => {
    const { user, validator } = setup();
    const labels = within(screen.getByRole('group', { name: 'Labels' }));
    await user.type(labels.getByPlaceholderText('label value'), 'some-value');
    const dbServiceLabelsGroup = within(
      screen.getByRole('group', { name: 'Database Service Labels' })
    );
    await user.type(
      dbServiceLabelsGroup.getByPlaceholderText('label value'),
      'some-value'
    );
    await selectEvent.create(screen.getByLabelText('Database Roles'), '*', {
      createOptionText: 'Database Role: *',
    });
    act(() => validator.validate());
    expect(
      labels.getByPlaceholderText('label key')
    ).toHaveAccessibleDescription('required');
    expect(
      dbServiceLabelsGroup.getByPlaceholderText('label key')
    ).toHaveAccessibleDescription('required');
    expect(
      screen.getByText('Wildcard is not allowed in database roles')
    ).toBeInTheDocument();
  });
});

describe('WindowsDesktopAccessSection', () => {
  const setup = () => {
    const onChange = jest.fn();
    let validator: Validator;
    render(
      <StatefulSection<WindowsDesktopAccess, ResourceAccessValidationResult>
        component={WindowsDesktopAccessSection}
        defaultValue={newResourceAccess('windows_desktop', defaultRoleVersion)}
        onChange={onChange}
        validatorRef={v => {
          validator = v;
        }}
        validate={validateResourceAccess}
      />
    );
    return { user: userEvent.setup(), onChange, validator };
  };

  test('editing', async () => {
    const { user, onChange } = setup();
    await user.type(screen.getByPlaceholderText('label key'), 'os');
    await user.type(screen.getByPlaceholderText('label value'), 'win-xp');
    await selectEvent.create(screen.getByLabelText('Logins'), 'julio', {
      createOptionText: 'Login: julio',
    });
    expect(onChange).toHaveBeenLastCalledWith({
      kind: 'windows_desktop',
      labels: [{ name: 'os', value: 'win-xp' }],
      logins: [
        expect.objectContaining({ value: '{{internal.windows_logins}}' }),
        expect.objectContaining({ label: 'julio', value: 'julio' }),
      ],
      hideValidationErrors: true,
    } as WindowsDesktopAccess);
  });

  test('validation', async () => {
    const { user, validator } = setup();
    await user.type(screen.getByPlaceholderText('label value'), 'some-value');
    act(() => validator.validate());
    expect(
      screen.getByPlaceholderText('label key')
    ).toHaveAccessibleDescription('required');
  });
});

describe('GitHubOrganizationAccessSection', () => {
  const setup = () => {
    const onChange = jest.fn();
    let validator: Validator;
    render(
      <StatefulSection<GitHubOrganizationAccess, ResourceAccessValidationResult>
        component={GitHubOrganizationAccessSection}
        defaultValue={newResourceAccess('git_server', defaultRoleVersion)}
        onChange={onChange}
        validatorRef={v => {
          validator = v;
        }}
        validate={validateResourceAccess}
      />
    );
    return { user: userEvent.setup(), onChange, validator };
  };

  test('editing', async () => {
    const { onChange } = setup();
    await selectEvent.create(
      screen.getByLabelText('Organization Names'),
      'illuminati',
      {
        createOptionText: 'Organization: illuminati',
      }
    );
    expect(onChange).toHaveBeenLastCalledWith({
      kind: 'git_server',
      organizations: [
        expect.objectContaining({ value: '{{internal.github_orgs}}' }),
        expect.objectContaining({ label: 'illuminati', value: 'illuminati' }),
      ],
      hideValidationErrors: true,
    } as GitHubOrganizationAccess);
  });
});

const reactSelectValueContainer = (input: HTMLInputElement) =>
  // eslint-disable-next-line testing-library/no-node-access
  input.closest('.react-select__value-container');
