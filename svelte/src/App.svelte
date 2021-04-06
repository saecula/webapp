<script>
	import { onMount, onDestroy } from "svelte";

	let messages = [];
	let message = {
		title: "",
		body: "",
	};

	const socket = new WebSocket("ws://localhost:4000/ws");
	socket.addEventListener("message", function ({ data }) {
		const message = JSON.parse(data);
		messages = [message, ...messages];
	});

	const fetchMain = async () => {
		const response = await fetch(`http://localhost:4000/`);
		const parsedResponse = await response.json();
		return parsedResponse;
	};

	onMount(async () => {
		const prevMessages = await fetchMain();
		messages = prevMessages;
	});

	const sendMessage = () => {
		if (message.body) {
			socket.send(JSON.stringify(message));
			message = {
				title: "",
				body: "",
			};
		}
	};

	onDestroy(() => {
		socket.disconnect();
	});
</script>

<style>
	h1 {
		color: #ff3e00;
		font-size: 4em;
		font-weight: 100;
	}
	h2 {
		color: #00a17e;
		font-size: 2em;
		font-weight: 100;
	}

	.container {
		display: flex;
		flex-direction: column;
		width: 80vw;
		margin: auto;
		text-align: center;
		background: url("sixpack.png") no-repeat bottom left;
		background-size: 20vh;
	}
	.message {
		width: 400px;
		margin: 10px auto;
		border: none;
		outline: white;
	}

	::placeholder {
		color: rgb(195, 195, 195);
		font-weight: 180;
	}
	#submit {
		height: 64px;
		color: #6d6d6d;
		background-color: #8affe6;
	}

	.post {
		list-style: none;
		color: rgb(80, 80, 80);
		margin: 20px auto;
		width: 200px;
		overflow: auto;
		padding: 50px;
		border-radius: 50%;
	}
</style>

<div class="container">
	<h1>its a webapp</h1>
	<h2>write me a message!!!!1omg pls</h2>

	<input
		placeholder="title here (if you want)"
		id="title"
		class="message"
		bind:value={message.title} />
	<input
		placeholder="you know you want to send me something"
		id="body"
		class="message input"
		bind:value={message.body} />

	<button class="message" id="submit" type="button" on:click={sendMessage}>
		send ur message
	</button>

	<ul class="posts">
		{#each messages as msg}
			<li class="post">
				<div><strong>{msg.title}</strong></div>
				<div>{msg.body}</div>
			</li>
		{/each}
	</ul>
</div>
