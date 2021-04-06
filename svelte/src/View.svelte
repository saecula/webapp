<!--<script>
    import { onMount, onDestroy } from "svelte";

    let socket = io("http://localhost:5000/ws");

    onMount(() => {
        socket.on("message", (data) => {
            messages = [...messages, data.text];
        });
    });

    onDestroy(() => {
        socket.disconnect();
    });

    const emptyMsg = {
        title: "",
        body: "",
    };
    let msg = emptyMsg;
    let messages = [];

    const sendMessage = () => {
        if (msg) {
            socket.emit("add-message", {
                ...msg,
                time: new Date().toLocaleTimeString(),
            });
            msg = emptyMsg;
        }
    };
</script><h3>Svelte Chat</h3>

<ul class="list-group">
    {#each messages as msg}
        <li class="list-group-item">
            <div>
                <strong>{msg.from}</strong><span class="time">{msg.time}</span>
            </div>
            <div>{msg.content}</div>
        </li>

        <div class="form-group">
            <input
                bind:this={msg}
                class="form-control"
                placeholder="Message"
                id="msg" />
        </div>
        <button
            type="button"
            on:click={sendMessage()}
            class="btn btn-primary pull-right">Send Message</button>
    {/each}
</ul> -->
