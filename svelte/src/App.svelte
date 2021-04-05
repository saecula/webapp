<script>
	const fetchMain = (async () => {
		const response = await fetch(`http://localhost:4000/`);
		return await response.json();
	})();
</script>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>

{#await fetchMain}
	<h1>loading...</h1>
{:then data}
	<main>
		<h1>hello {data.join(' ')}!</h1>
		{#each data as title}
			<h2>{title}</h2>
		{/each}
		<p>this is a much nicer frontend!</p>
	</main>
{:catch error}
	<p>ohnoes!! {error}</p>
{/await}
