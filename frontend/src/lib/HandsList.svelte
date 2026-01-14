<script lang="ts">
  import { onMount } from 'svelte';
  import { GetHands, GetHandByID } from '../../wailsjs/go/main/App';

  let hands: any[] = [];
  let loading = true;
  let selectedHand: any = null;
  let filter = {
    limit: 50,
    offset: 0
  };

  onMount(async () => {
    await loadHands();
  });

  async function loadHands() {
    loading = true;
    try {
      hands = await GetHands(filter);
    } catch (err) {
      console.error('Error loading hands:', err);
    } finally {
      loading = false;
    }
  }

  async function selectHand(id: number) {
    try {
      selectedHand = await GetHandByID(id);
    } catch (err) {
      console.error('Error loading hand details:', err);
    }
  }

  function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleString();
  }

  function formatAmount(amount: number): string {
    const sign = amount >= 0 ? '+' : '';
    return `${sign}$${amount.toFixed(2)}`;
  }

  function getResultClass(amount: number): string {
    return amount >= 0 ? 'text-success-500' : 'text-error-500';
  }
</script>

<div class="container mx-auto p-6 h-full overflow-auto">
  <div class="mb-4">
    <h2 class="text-xl font-bold mb-4">Hand History</h2>
    <button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700" on:click={loadHands}>
      Refresh
    </button>
  </div>

  {#if loading}
    <div class="text-center py-8">
      <p class="text-surface-400">Loading hands...</p>
    </div>
  {:else if hands.length === 0}
    <div class="bg-gray-800 rounded-lg p-8 text-center">
      <p class="text-gray-400 mb-4">No hands found</p>
      <p class="text-sm text-gray-500">
        Make sure your hand history directories are configured in Settings
        and that you've played some hands.
      </p>
    </div>
  {:else}
    <div class="overflow-x-auto bg-gray-800 rounded-lg">
      <table class="w-full text-left text-white">
        <thead class="bg-gray-700">
          <tr>
            <th class="px-4 py-3">Date/Time</th>
            <th class="px-4 py-3">Hand ID</th>
            <th class="px-4 py-3">Game Type</th>
            <th class="px-4 py-3">Stakes</th>
            <th class="px-4 py-3">Table</th>
            <th class="px-4 py-3">Hero</th>
            <th class="px-4 py-3">Result</th>
          </tr>
        </thead>
        <tbody>
          {#each hands as hand}
            <tr on:click={() => selectHand(hand.id)} class="cursor-pointer hover:bg-gray-700 border-b border-gray-700">
              <td class="px-4 py-3">{formatDate(hand.date_time)}</td>
              <td class="px-4 py-3 font-mono text-sm">{hand.hand_id}</td>
              <td class="px-4 py-3">{hand.game_type || 'N/A'}</td>
              <td class="px-4 py-3">{hand.stakes || 'N/A'}</td>
              <td class="px-4 py-3">{hand.table_name || 'N/A'}</td>
              <td class="px-4 py-3">{hand.hero_name || 'N/A'}</td>
              <td class="px-4 py-3 {getResultClass(hand.result)}">
                {formatAmount(hand.result)}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}

  <!-- Hand Detail Modal -->
  {#if selectedHand}
    <div class="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50" on:click={() => selectedHand = null}>
      <div class="bg-gray-800 rounded-lg p-6 max-w-3xl max-h-[90vh] overflow-y-auto" on:click|stopPropagation>
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-xl font-bold text-white">Hand #{selectedHand.hand_id}</h3>
          <button class="text-gray-400 hover:text-white text-2xl" on:click={() => selectedHand = null}>
            ✕
          </button>
        </div>
        <div class="space-y-4 text-white">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <p class="text-sm text-gray-400">Game Type</p>
              <p class="font-semibold">{selectedHand.game_type}</p>
            </div>
            <div>
              <p class="text-sm text-gray-400">Stakes</p>
              <p class="font-semibold">{selectedHand.stakes}</p>
            </div>
            <div>
              <p class="text-sm text-gray-400">Table</p>
              <p class="font-semibold">{selectedHand.table_name}</p>
            </div>
            <div>
              <p class="text-sm text-gray-400">Result</p>
              <p class="font-semibold {getResultClass(selectedHand.result)}">
                {formatAmount(selectedHand.result)}
              </p>
            </div>
          </div>
          
          {#if selectedHand.hole_cards}
            <div>
              <p class="text-sm text-gray-400 mb-2">Hole Cards</p>
              <div class="flex gap-2">
                {#each JSON.parse(selectedHand.hole_cards) as card}
                  <span class="px-3 py-1 bg-blue-600 rounded">{card}</span>
                {/each}
              </div>
            </div>
          {/if}
          
          {#if selectedHand.board}
            <div>
              <p class="text-sm text-gray-400 mb-2">Board</p>
              <div class="flex gap-2">
                {#each JSON.parse(selectedHand.board) as card}
                  <span class="px-3 py-1 bg-purple-600 rounded">{card}</span>
                {/each}
              </div>
            </div>
          {/if}

          <div>
            <p class="text-sm text-gray-400 mb-2">Raw Hand History</p>
            <pre class="bg-gray-900 p-4 rounded overflow-auto max-h-64 text-xs">{selectedHand.raw_text}</pre>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>
