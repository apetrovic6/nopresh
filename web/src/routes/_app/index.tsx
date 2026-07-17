import { createFileRoute } from '@tanstack/react-router'
import { useSuspenseQuery } from "@connectrpc/connect-query";
import { me } from '#/gen/proto/auth/v1/auth-AuthService_connectquery';
import { getSettings } from '#/gen/proto/settings/v1/settings-SettingsService_connectquery';

export const Route = createFileRoute('/_app/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { data: currentUser } = useSuspenseQuery(me, {}, { staleTime: 5 * 60 * 1000 });
  const { data: userSettings} = useSuspenseQuery(getSettings, {}, { staleTime: 5 * 60 * 1000 });

  return (
    <div className='space-y-4 flex-1 min-h-0 overflow-auto p-2'>
      <h1>WIP</h1>
      <div>
        <div>
          {currentUser.email}
          {currentUser.name}
        </div>

      <div>
        {userSettings.settings?.defaultMedicationId}
        </div>
      </div>
    </div>
  )
}
