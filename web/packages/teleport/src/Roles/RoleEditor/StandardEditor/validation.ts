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

import { Option } from 'shared/components/Select';
import {
  arrayOf,
  requiredField,
  RuleSetValidationResult,
  runRules,
  ValidationResult,
} from 'shared/components/Validation/rules';

import { nonEmptyLabels } from 'teleport/components/LabelsInput/LabelsInput';
import { KubernetesResourceKind } from 'teleport/services/resources';

import {
  KubernetesResourceModel,
  MetadataModel,
  ResourceAccess,
  RoleEditorModel,
  RuleModel,
} from './standardmodel';

const kubernetesClusterWideResourceKinds: KubernetesResourceKind[] = [
  'namespace',
  'kube_node',
  'persistentvolume',
  'clusterrole',
  'clusterrolebinding',
  'certificatesigningrequest',
];

export type RoleEditorModelValidationResult = {
  metadata: MetadataValidationResult;
  resources: ResourceAccessValidationResult[];
  rules: AccessRuleValidationResult[];
};

/**
 * Validates the role editor model. In addition to the model itself, this
 * function also takes the previous model and previous validation result. The
 * intention here is to only return a newly created result if the previous
 * model is indeed different. This pattern is then repeated in other validation
 * functions.
 *
 * The purpose of this is less about the performance of the validation process
 * itself, and more about enabling memoization-based rendering optimizations:
 * UI components that take either entire or partial validation results can be
 * cached if the validation results don't change.
 *
 * Note that we can't use `useMemo` here, because `validateRoleEditorModel` is
 * called from the state reducer. Also `highbar.memoize` was not applicable, as
 * it caches an arbitrary amount of previous results.
 */
export function validateRoleEditorModel(
  model: RoleEditorModel,
  previousModel: RoleEditorModel | undefined,
  previousResult: RoleEditorModelValidationResult | undefined
): RoleEditorModelValidationResult {
  return {
    metadata: validateMetadata(
      model.metadata,
      previousModel?.metadata,
      previousResult?.metadata
    ),
    resources: validateResourceAccessList(
      model.resources,
      previousModel?.resources,
      previousResult?.resources
    ),
    rules: validateAccessRuleList(
      model.rules,
      previousModel?.rules,
      previousResult?.rules
    ),
  };
}

function validateMetadata(
  model: MetadataModel,
  previousModel: MetadataModel,
  previousResult: MetadataValidationResult
): MetadataValidationResult {
  if (previousModel === model) {
    return previousResult;
  }
  return runRules(model, metadataRules);
}

const metadataRules = {
  name: requiredField('Role name is required'),
  labels: nonEmptyLabels,
};
export type MetadataValidationResult = RuleSetValidationResult<
  typeof metadataRules
>;

function validateResourceAccessList(
  resources: ResourceAccess[],
  previousResources: ResourceAccess[],
  previousResults: ResourceAccessValidationResult[]
): ResourceAccessValidationResult[] {
  if (previousResources === resources) {
    return previousResults;
  }
  return resources.map((res, i) =>
    validateResourceAccess(res, previousResources?.[i], previousResults?.[i])
  );
}

export function validateResourceAccess(
  resource: ResourceAccess,
  previousResource: ResourceAccess,
  previousResult: ResourceAccessValidationResult
): ResourceAccessValidationResult {
  if (resource === previousResource) {
    return previousResult;
  }

  const { kind } = resource;
  switch (kind) {
    case 'kube_cluster':
      return runRules(resource, kubernetesAccessValidationRules);
    case 'node':
      return runRules(resource, serverAccessValidationRules);
    case 'app':
      return runRules(resource, appAccessValidationRules);
    case 'db':
      return runRules(resource, databaseAccessValidationRules);
    case 'windows_desktop':
      return runRules(resource, windowsDesktopAccessValidationRules);
    default:
      kind satisfies never;
  }
}

export type ResourceAccessValidationResult =
  | ServerAccessValidationResult
  | KubernetesAccessValidationResult
  | AppAccessValidationResult
  | DatabaseAccessValidationResult
  | WindowsDesktopAccessValidationResult;

const validKubernetesResource = (res: KubernetesResourceModel) => () => {
  const name = requiredField(
    'Resource name is required, use "*" for any resource'
  )(res.name)();
  const namespace = kubernetesClusterWideResourceKinds.includes(res.kind.value)
    ? { valid: true }
    : requiredField('Namespace is required for resources of this kind')(
        res.namespace
      )();
  return {
    valid: name.valid && namespace.valid,
    name,
    namespace,
  };
};
export type KubernetesResourceValidationResult = {
  name: ValidationResult;
  namespace: ValidationResult;
};

const kubernetesAccessValidationRules = {
  labels: nonEmptyLabels,
  resources: arrayOf(validKubernetesResource),
};
export type KubernetesAccessValidationResult = RuleSetValidationResult<
  typeof kubernetesAccessValidationRules
>;

const noWildcard = (message: string) => (value: string) => () => {
  const valid = value !== '*';
  return { valid, message: valid ? '' : message };
};

const noWildcardOptions = (message: string) => (options: Option[]) => () => {
  const valid = options.every(o => o.value !== '*');
  return { valid, message: valid ? '' : message };
};

const serverAccessValidationRules = {
  labels: nonEmptyLabels,
  logins: noWildcardOptions('Wildcard is not allowed in logins'),
};
export type ServerAccessValidationResult = RuleSetValidationResult<
  typeof serverAccessValidationRules
>;

const appAccessValidationRules = {
  labels: nonEmptyLabels,
  awsRoleARNs: arrayOf(noWildcard('Wildcard is not allowed in AWS role ARNs')),
  azureIdentities: arrayOf(
    noWildcard('Wildcard is not allowed in Azure identities')
  ),
  gcpServiceAccounts: arrayOf(
    noWildcard('Wildcard is not allowed in GCP service accounts')
  ),
};
export type AppAccessValidationResult = RuleSetValidationResult<
  typeof appAccessValidationRules
>;

const databaseAccessValidationRules = {
  labels: nonEmptyLabels,
  roles: noWildcardOptions('Wildcard is not allowed in database roles'),
};
export type DatabaseAccessValidationResult = RuleSetValidationResult<
  typeof databaseAccessValidationRules
>;

const windowsDesktopAccessValidationRules = {
  labels: nonEmptyLabels,
};
export type WindowsDesktopAccessValidationResult = RuleSetValidationResult<
  typeof windowsDesktopAccessValidationRules
>;

export function validateAccessRuleList(
  rules: RuleModel[],
  previousRules: RuleModel[],
  previousResults: AccessRuleValidationResult[]
): AccessRuleValidationResult[] {
  if (previousRules === rules) {
    return previousResults;
  }
  return rules.map((rule, i) =>
    validateAccessRule(rule, previousRules?.[i], previousResults?.[i])
  );
}

export const validateAccessRule = (
  rule: RuleModel,
  previousRule: RuleModel,
  previousResult: AccessRuleValidationResult
): AccessRuleValidationResult => {
  if (previousRule === rule) {
    return previousResult;
  }
  return runRules(rule, accessRuleValidationRules);
};

const accessRuleValidationRules = {
  resources: requiredField('At least one resource kind is required'),
  verbs: requiredField('At least one permission is required'),
};
export type AccessRuleValidationResult = RuleSetValidationResult<
  typeof accessRuleValidationRules
>;
