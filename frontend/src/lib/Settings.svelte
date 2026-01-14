<script lang="ts">
  import { onMount } from 'svelte';
  import { GetConfig, UpdateConfig, SelectDirectory, GetWatcherStatus } from '../../wailsjs/go/main/App';

  let config: any = null;
  let watcherStatus: any = null;
  let loading = true;
  let saving = false;

  onMount(async () => {
    await loadConfig();
    await loadWatcherStatus();
  });

  async function loadConfig() {
    loading = true;
    try {
      config = await GetConfig();
    } catch (err) {
      console.error('Error loading config:', err);
    } finally {
      loading = false;
    }
  }

  async function loadWatcherStatus() {
    try {
      watcherStatus = await GetWatcherStatus();
    } catch (err) {
      console.error('Error loading watcher status:', err);
    }
  }

  async function saveConfig() {
    saving = true;
    try {
      await UpdateConfig(config);
      await loadWatcherStatus();
      alert('Settings saved successfully!');
    } catch (err) {
      console.error('Error saving config:', err);
      alert('Failed to save settings: ' + err);
    } finally {
      saving = false;
    }
  }

  async function selectDirectory(siteName: string) {
    try {
      const path = await SelectDirectory();
      if (path) {
        config.sites[siteName].watch_path = path;
      }
    } catch (err) {
      console.error('Error selecting directory:', err);
    }
  }

  setInterval(loadWatcherStatus, 5000); // Update status every 5 seconds
</script>

<div class="container mx-auto p-6 h-full overflow-auto text-white">
  <h2 class="text-xl font-bold mb-6">Settings</h2>

  {#if loading}
    <div class="text-center py-8">
      <p class="text-gray-400">Loading settings...</p>
    </div>
  {:else if config}
    <div class="space-y-6">
      <!-- General Settings -->
      <div class="bg-gray-800 rounded-lg p-6">
        <h3 class="text-lg font-semibold mb-4">General</h3>
        
        <label class="block mb-4">
          <span class="block mb-2">Hero Name (Default)</span>
          <input
            class="w-full px-4 py-2 bg-gray-700 text-white rounded border border-gray-600 focus:border-blue-500 focus:outline-none"
            type="text"
            placeholder="Your poker screen name..."
            bind:value={config.hero_name}
          />
          <p class="text-sm text-gray-400 mt-1">
            Your primary screen name used to identify you in hand histories
          </p>
        </label>

        <label class="block">
          <span class="block mb-2">Theme</span>
          <select class="w-full px-4 py-2 bg-gray-700 text-white rounded border border-gray-600 focus:border-blue-500 focus:outline-none" bind:value={config.theme}>
            <option value="dark">Dark</option>
            <option value="light">Light</option>
          </select>
        </label>
      </div>

      <!-- Poker Sites Configuration -->
      <div class="bg-gray-800 rounded-lg p-6">
        <h3 class="text-lg font-semibold mb-4">Poker Sites</h3>

        {#each Object.entries(config.sites) as [siteName, site]}
          <div class="bg-gray-700 bg-opacity-50 rounded-lg p-4 mb-4">
            <div class="flex items-center justify-between mb-3">
              <h4 class="font-semibold">{site.name}</h4>
              <label class="flex items-center space-x-2">
                <input
                  class="w-4 h-4"
                  type="checkbox"
                  bind:checked={site.enabled}
                />
                <span>Enabled</span>
              </label>
            </div>

            <label class="block">
              <span class="text-sm block mb-2">Hand History Directory</span>
              <div class="flex gap-2">
                <input
                  class="flex-1 px-4 py-2 bg-gray-700 text-white rounded border border-gray-600 focus:border-blue-500 focus:outline-none disabled:bg-gray-800 disabled:text-gray-500"
                  type="text"
                  placeholder="Select directory..."
                  bind:value={site.watch_path}
                  disabled={!site.enabled}
                />
                <button
                  class="px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700 disabled:bg-gray-600 disabled:cursor-not-allowed"
                  on:click={() => selectDirectory(siteName)}
                  disabled={!site.enabled}
                >
                  Browse
                </button>
              </div>
            </label>
          </div>
        {/each}
      </div>

      <!-- Watcher Status -->
      {#if watcherStatus}
        <div class="bg-gray-800 rounded-lg p-6">
          <h3 class="text-lg font-semibold mb-4">File Watcher Status</h3>
          
          <div class="space-y-2">
            <div class="flex items-center gap-2">
              <span class="text-sm text-gray-400">Status:</span>
              <span class="px-3 py-1 rounded text-sm {watcherStatus.is_running ? 'bg-green-600' : 'bg-red-600'}">
                {watcherStatus.is_running ? 'Running' : 'Stopped'}
              </span>
            </div>

            <div>
              <span class="text-sm text-gray-400">Queue Length:</span>
              <span class="font-semibold ml-2">{watcherStatus.queue_length}</span>
            </div>

            {#if watcherStatus.watched_paths && Object.keys(watcherStatus.watched_paths).length > 0}
              <div>
                <p class="text-sm text-gray-400 mb-2">Watching:</p>
                <ul class="list-disc list-inside">
                  {#each Object.keys(watcherStatus.watched_paths) as path}
                    <li class="text-sm">{path}</li>
                  {/each}
                </ul>
              </div>
            {/if}
          </div>
        </div>
      {/if}

      <!-- Database Info -->
      <div class="bg-gray-800 rounded-lg p-6">
        <h3 class="text-lg font-semibold mb-4">Database</h3>
        <div class="space-y-2">
          <div>
            <span class="text-sm text-gray-400">Database Path:</span>
            <p class="font-mono text-sm mt-1 break-all">{config.database_path}</p>
          </div>
        </div>
      </div>

      <!-- Save Button -->
      <div class="flex justify-end">
        <button
          class="px-6 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed"
          on:click={saveConfig}
          disabled={saving}
        >
          {saving ? 'Saving...' : 'Save Settings'}
        </button>
      </div>
    </div>
  {/if}
</div>
