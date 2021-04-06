<script>
	import { onMount, onDestroy } from "svelte";
	const uuidv4 = new RegExp(
		/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i
	);

	let socket;
	let messages = [];
	let message = {
		title: "",
		body: "",
	};

	onMount(async () => {
		initSocket();
		const res = await fetchMain();
		const formatted = formatRes(res);
		messages = formatted;
	});

	onDestroy(() => {
		socket.disconnect();
	});

	const initSocket = () => {
		socket = new WebSocket("ws://localhost:4000/ws");
		socket.addEventListener("message", function ({ data }) {
			const message = JSON.parse(data);
			message.title = unwrap(message.title);
			messages = [message, ...messages];
		});
	};

	const fetchMain = async () => {
		const response = await fetch(`http://localhost:4000/`);
		const parsedResponse = await response.json();
		return parsedResponse;
	};

	const sendMessage = () => {
		const { title, body } = message;
		if (title) {
			message.title = wrap(title);
		}
		if (body) {
			socket.send(JSON.stringify(message));
			message = {
				title: "",
				body: "",
			};
		}
	};

	const formatRes = (strArr) =>
		strArr.map(({ title, ...rest }) => ({
			title: unwrap(title),
			...rest,
		}));
	const wrap = (str) => str.replaceAll(" ", "-").trim();
	const unwrap = (str) => (str.match(uuidv4) ? "" : str.replaceAll("-", " "));
</script>

<style>
	.credit {
		color: #989898;
		font-weight: 150;
		font-size: 0.5rem;
		margin-left: 83%;
	}
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
		min-height: 100vh;
		margin: auto;
		text-align: center;
		background: url("sixpack.png") no-repeat left calc(100vh - 231px);
		background-size: 200px;
	}
	.message {
		width: 400px;
		max-width: 80vw;
		margin: 10px auto;
		border: none;
		outline: white;
		font-weight: 180;
	}
	::placeholder {
		color: rgb(195, 195, 195);
	}
	.submit {
		height: 64px;
		color: #626262;
		background-color: #8affe6;
		font-weight: 300;
	}

	button:disabled,
	button[disabled] {
		color: #bbbbbb;
	}

	.posts {
		margin: auto;
		padding: 0;
	}

	.post {
		list-style: none;
		color: rgb(80, 80, 80);
		background-color: #ffffff61;
		margin: 20px auto;
		width: 200px;
		overflow: auto;
		padding: 50px;
		border-radius: 50%;
	}
</style>

<div class="credit">
	{'ripped gopher by'}
	<a
		id="link"
		href="https://www.redbubble.com/people/clgtart/shop">clgtart</a>
</div>
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

	<button
		class="message submit"
		type="button"
		disabled={message.body === ''}
		on:click={sendMessage}>
		send itt
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
