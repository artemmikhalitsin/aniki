<script lang="ts">
  import { onMount } from 'svelte';
  import { GetStats } from '../../wailsjs/go/main/App';

  let heroName = '';
  let stats: any = null;
  let loading = false;

  async function loadStats() {
    if (!heroName) return;

    loading = true;
    try {
      stats = await GetStats(heroName);
    } catch (err) {
      console.error('Error loading stats:', err);
    } finally {
      loading = false;
    }
  }

  function formatAmount(amount: number): string {
    const sign = amount >= 0 ? '+' : '';
    return `${sign}$${amount.toFixed(2)}`;
  }
</script>

<div class="container mx-auto p-6">
  <h2 class="text-xl font-bold mb-6 text-white">Player Statistics</h2>

  <div class="bg-gray-800 rounded-lg p-6 mb-6">
    <label class="block mb-2">
      <span class="text-white mb-2 block">Hero Name</span>
      <div class="flex gap-2">
        <input
          class="flex-1 px-4 py-2 bg-gray-700 text-white rounded border border-gray-600 focus:border-blue-500 focus:outline-none"
          type="text"
          placeholder="Enter player name..."
          bind:value={heroName}
          on:keypress={(e) => e.key === 'Enter' && loadStats()}
        />
        <button
          class="px-6 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed"
          on:click={loadStats}
          disabled={!heroName || loading}
        >
          {loading ? 'Loading...' : 'Load Stats'}
        </button>
      </div>
    </label>
  </div>

  {#if stats}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div class="bg-gray-800 rounded-lg p-6">
        <p class="text-sm text-gray-400 mb-2">Total Hands</p>
        <p class="text-3xl font-bold text-white">{stats.total_hands}</p>
      </div>

      <div class="bg-gray-800 rounded-lg p-6">
        <p class="text-sm text-gray-400 mb-2">Total Won/Lost</p>
        <p class="text-3xl font-bold {stats.total_won >= 0 ? 'text-green-400' : 'text-red-400'}">
          {formatAmount(stats.total_won)}
        </p>
      </div>

      <div class="bg-gray-800 rounded-lg p-6">
        <p class="text-sm text-gray-400 mb-2">Total Rake Paid</p>
        <p class="text-3xl font-bold text-yellow-400">
          ${stats.total_rake.toFixed(2)}
        </p>
      </div>

      <div class="bg-gray-800 rounded-lg p-6">
        <p class="text-sm text-gray-400 mb-2">Biggest Win</p>
        <p class="text-3xl font-bold text-green-400">
          {formatAmount(stats.biggest_win)}
        </p>
      </div>

      <div class="bg-gray-800 rounded-lg p-6">
        <p class="text-sm text-gray-400 mb-2">Biggest Loss</p>
        <p class="text-3xl font-bold text-red-400">
          {formatAmount(stats.biggest_loss)}
        </p>
      </div>

      <div class="bg-gray-800 rounded-lg p-6">
        <p class="text-sm text-gray-400 mb-2">Win Rate</p>
        <p class="text-3xl font-bold text-white">
          {stats.win_rate.toFixed(2)}%
        </p>
      </div>

      <div class="bg-gray-800 rounded-lg p-6">
        <p class="text-sm text-gray-400 mb-2">Hands Won</p>
        <p class="text-3xl font-bold text-green-400">
          {stats.hands_won}
        </p>
      </div>

      <div class="bg-gray-800 rounded-lg p-6">
        <p class="text-sm text-gray-400 mb-2">Hands Lost</p>
        <p class="text-3xl font-bold text-red-400">
          {stats.hands_lost}
        </p>
      </div>

      <div class="bg-gray-800 rounded-lg p-6">
        <p class="text-sm text-gray-400 mb-2">Win Percentage</p>
        <p class="text-3xl font-bold text-white">
          {stats.total_hands > 0 ? ((stats.hands_won / stats.total_hands) * 100).toFixed(1) : 0}%
        </p>
      </div>
    </div>
  {:else}
    <div class="bg-gray-800 rounded-lg p-8 text-center">
      <p class="text-gray-400">
        Enter a hero name above to view statistics
      </p>
    </div>
  {/if}
</div>
