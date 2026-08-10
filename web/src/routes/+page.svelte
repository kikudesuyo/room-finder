<script>
	import { onMount } from 'svelte';

	const apiBaseURL = import.meta.env.PUBLIC_API_BASE_URL || 'http://localhost:8081';
	let profiles = $state([]);
	let offers = $state([]);
	let selectedProfileID = $state(null);
	let loading = $state(true);
	let errorMessage = $state('');

	async function fetchJSON(path) {
		const response = await globalThis.fetch(`${apiBaseURL}${path}`);
		if (!response.ok) {
			throw new Error(`API request failed: ${response.status}`);
		}
		const body = await response.json();
		return body.data;
	}

	async function loadOffers(profileID) {
		selectedProfileID = profileID;
		offers = await fetchJSON(`/api/v1/search-profiles/${profileID}/rental-offers`);
	}

	async function load() {
		loading = true;
		errorMessage = '';
		try {
			profiles = await fetchJSON('/api/v1/search-profiles');
			if (profiles.length > 0) {
				await loadOffers(profiles[0].id);
			} else {
				offers = [];
				selectedProfileID = null;
			}
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'データを取得できませんでした';
		} finally {
			loading = false;
		}
	}

	function yen(value) {
		return value === null ? '—' : `${value.toLocaleString('ja-JP')}円`;
	}

	onMount(load);
</script>

<svelte:head>
	<title>Room Finder | 賃貸募集</title>
	<meta name="description" content="検索条件に一致した賃貸募集情報を確認する" />
</svelte:head>

<div class="mx-auto min-h-screen max-w-6xl px-5 py-8 sm:px-8 sm:py-12">
	<header class="border-on-surface/10 mb-10 flex items-end justify-between gap-4 border-b pb-6">
		<div>
			<p class="font-label-md text-primary tracking-[0.18em] uppercase">Room Finder</p>
			<h1 class="font-headline-lg mt-3">条件に合う賃貸募集</h1>
			<p class="text-body-md text-on-surface-variant mt-2">Agentが取得・判定した現在の情報です。</p>
		</div>
		<button class="button-secondary" type="button" onclick={load} disabled={loading}
			>再読み込み</button
		>
	</header>

	{#if errorMessage}
		<div class="notice-error" role="alert">{errorMessage}</div>
	{:else if loading}
		<p class="text-on-surface-variant">読み込み中…</p>
	{:else}
		<div class="grid gap-8 lg:grid-cols-[18rem_1fr]">
			<aside>
				<div class="mb-3 flex items-center justify-between">
					<h2 class="font-label-lg">検索条件</h2>
					<span class="text-on-surface-variant text-sm">{profiles.length}件</span>
				</div>
				{#if profiles.length === 0}
					<div class="empty-card">検索条件がまだありません。</div>
				{:else}
					<div class="space-y-2">
						{#each profiles as profile (profile.id)}
							<button
								class:selected={selectedProfileID === profile.id}
								class="profile-card"
								type="button"
								onclick={() => loadOffers(profile.id)}
							>
								<span class="text-primary text-sm">条件 #{profile.id}</span>
								<span class="text-on-surface-variant mt-1 line-clamp-3 text-left text-sm"
									>{profile.initial_prompt}</span
								>
							</button>
						{/each}
					</div>
				{/if}
			</aside>

			<section aria-live="polite">
				<div class="mb-3 flex items-center justify-between">
					<h2 class="font-label-lg">一致した募集</h2>
					<span class="text-on-surface-variant text-sm">{offers.length}件</span>
				</div>
				{#if offers.length === 0}
					<div class="empty-card">条件に一致する募集はまだありません。</div>
				{:else}
					<div class="grid gap-4 md:grid-cols-2">
						{#each offers as offer (offer.id)}
							<article class="offer-card">
								<div class="flex items-start justify-between gap-4">
									<div>
										<h3 class="font-label-lg">{offer.name || '名称未取得'}</h3>
										<p class="text-on-surface-variant mt-1 text-sm">
											{offer.address || '住所未取得'}
										</p>
									</div>
									<span class="tag">{offer.source}</span>
								</div>
								<dl class="offer-facts">
									<div>
										<dt>家賃</dt>
										<dd>{yen(offer.rent_yen)}</dd>
									</div>
									<div>
										<dt>管理費</dt>
										<dd>{yen(offer.management_fee_yen)}</dd>
									</div>
									<div>
										<dt>間取り</dt>
										<dd>{offer.room_layout || '—'}</dd>
									</div>
									<div>
										<dt>面積</dt>
										<dd>
											{offer.area_square_meters === null ? '—' : `${offer.area_square_meters}㎡`}
										</dd>
									</div>
								</dl>
								<button
									class="source-link"
									type="button"
									onclick={() => globalThis.open(offer.source_url, '_blank', 'noopener,noreferrer')}
									>元サイトで詳細を見る ↗</button
								>
								<p class="text-on-surface-variant mt-3 text-xs">
									取得 {new Date(offer.captured_at).toLocaleString('ja-JP')}
								</p>
							</article>
						{/each}
					</div>
				{/if}
			</section>
		</div>
	{/if}
</div>
