/**
 * Teleport
 * Copyright (C) 2023  Gravitational, Inc.
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

import { makeEmptyAttempt } from 'shared/hooks/useAsync';

import { UserContextValue } from 'teleport/User/UserContext';

export const makeTestUserContext = (
  overrides?: Partial<UserContextValue>
): UserContextValue => {
  return Object.assign(
    {
      preferences: {
        theme: 1,
        onboard: {
          preferredResources: [],
        },
        keyboardLayout: 0,
      },
      clusterPreferences: {
        pinnedResources: [],
      },
      updateClusterPreferencesAttempt: makeEmptyAttempt(),
      updatePreferences: () => Promise.resolve(),
      updateClusterPinnedResources: () => Promise.resolve(),
      getClusterPinnedResources: () => Promise.resolve(),
      updateDiscoverResourcePreferences: () => Promise.resolve(),
    },
    overrides
  );
};
