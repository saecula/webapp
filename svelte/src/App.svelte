<script>
	import { onMount, onDestroy } from "svelte";
	// import io from "socket.io-client";

	const socket = new WebSocket("ws://localhost:4000/ws");
	socket.addEventListener("open", function (event) {
		console.log("omg opened! what is event:", event);
	});
	socket.addEventListener("message", function ({ data }) {
		console.log("omg NEWWW incoming message:", data);
		const message = JSON.parse(data);
		messages = [...messages, message];
	});

	let messages = [];
	let message = {
		title: "",
		body: "",
	};

	const fetchMain = async () => {
		const response = await fetch(`http://localhost:4000/`);
		console.log("response...", response);
		const parsed = await response.json();
		console.log("parsed", parsed);
		return parsed;
	};

	onMount(async () => {
		const prevMessages = await fetchMain();
		console.log("ooo", prevMessages);
		messages = prevMessages;
	});

	onDestroy(() => {
		socket.disconnect();
	});

	const sendMessage = () => {
		console.log("wutttt", message);
		if (message.body) {
			socket.send(JSON.stringify(message));
		}
	};
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
</style>

<h1>its a webapp</h1>
<h2>write a message!!!!omg</h2>

<div class="form-group">
	<input
		class="form-control"
		placeholder="TITLE"
		id="body"
		bind:value={message.title} />
	<input
		class="form-control"
		placeholder="you know you want to"
		id="body"
		bind:value={message.body} />
</div>

<button type="button" on:click={sendMessage}> send ur message </button>

<ul class="list-group">
	{#each messages as msg}
		<li class="list-group-item">
			<div><strong>{msg.title}</strong></div>
			<div>{msg.body}</div>
		</li>
	{/each}
</ul>
<!-- {#await fetchMain}
	<h1>loading...</h1>
{:then data}
	<main>
		<h1>hello!</h1>
		{#each data as title}
			<h2>{title}</h2>
		{/each}
		<p>this is a much nicer frontend!</p>
	</main>
{:catch error}
	<p>ohnoes!! {error}</p>
{/await} -->
