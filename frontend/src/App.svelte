<script>
  // Greet ではなく、新しく作った GetTasks と AddTask をインポート
  import { GetTasks, AddTask } from '../wailsjs/go/main/App.js'
  import { onMount } from 'svelte';

  let tasks = [];
  let newTaskContent = "";

  // アプリ起動時にタスクを読み込む
  onMount(async () => {
    await loadTasks();
  });

  // タスク取得処理
  async function loadTasks() {
    tasks = await GetTasks();
  }

  // 追加ボタンを押した時の処理
  async function addTask() {
    if (!newTaskContent) return;
    // Goに追加してもらって、最新リストをもらう
    tasks = await AddTask(newTaskContent);
    newTaskContent = ""; // 入力欄を空にする
  }
</script>

<main>
  <img alt="Wails logo" src="./assets/images/logo-universal.png" class="logo" />

  <h1>notree Task Manager</h1>

  <div class="input-box">
    <input
      bind:value={newTaskContent}
      placeholder="学習タスクを入力..."
      type="text"
    />
    <button on:click={addTask}>追加</button>
  </div>

  <ul class="task-list">
    {#each tasks as task}
      <li>
        <input type="checkbox" checked={task.is_done} />
        <span>{task.content}</span>
      </li>
    {/each}
  </ul>
</main>

<style>
  .logo { display: block; width: 100px; margin: 0 auto; }
  main { padding: 20px; max-width: 600px; margin: 0 auto; font-family: sans-serif; }
  h1 { text-align: center; }
  .input-box { display: flex; gap: 10px; margin-bottom: 20px; }
  input[type="text"] { flex: 1; padding: 10px; font-size: 16px; }
  button { padding: 10px 20px; cursor: pointer; background-color: #3e3e3e; color: white; border: none; border-radius: 4px;}
  .task-list { list-style: none; padding: 0; }
  .task-list li { background: #f9f9f9; padding: 10px; margin-bottom: 5px; border-radius: 4px; display: flex; align-items: center; gap: 10px; }
</style>