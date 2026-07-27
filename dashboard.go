// dashboard.go — embedded HTML dashboard (GitHub-dark style).
package main

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>GLM-ZAI-2API</title>
<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #0d1117; color: #c9d1d9; min-height: 100vh; padding: 24px; }
.header { display: flex; align-items: center; gap: 12px; margin-bottom: 24px; }
.header h1 { font-size: 20px; font-weight: 600; color: #e6edf3; }
.header .sub { font-size: 13px; color: #8b949e; }
.dot { width: 10px; height: 10px; border-radius: 50%; background: #56d364; box-shadow: 0 0 8px #56d364; animation: pulse 2s infinite; }
@keyframes pulse { 0%,100%{opacity:1} 50%{opacity:.5} }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 16px; }
.card { background: #161b22; border: 1px solid #21262d; border-radius: 8px; padding: 18px; }
.card-title { font-size: 11px; font-weight: 600; text-transform: uppercase; letter-spacing: .05em; color: #8b949e; margin-bottom: 14px; }
.row { display: flex; justify-content: space-between; align-items: center; padding: 6px 0; border-bottom: 1px solid #21262d; }
.row:last-child { border-bottom: none; }
.row-label { font-size: 13px; color: #8b949e; }
.row-value { font-size: 13px; color: #e6edf3; font-weight: 500; }
.badge { padding: 2px 8px; border-radius: 12px; font-size: 11px; font-weight: 600; cursor: pointer; transition: opacity .15s; user-select: none; }
.badge:hover { opacity: .75; }
.badge-green  { background: #1a3a1a; color: #56d364; border: 1px solid #238636; }
.badge-red    { background: #3a1a1a; color: #f85149; border: 1px solid #da3633; }
.badge-gray   { background: #21262d; color: #8b949e; border: 1px solid #30363d; }
.badge-blue   { background: #1a2a3a; color: #79c0ff; border: 1px solid #1f6feb; }
.badge-yellow { background: #2a2000; color: #e3b341; border: 1px solid #9e6a03; }
.endpoint { font-family: monospace; font-size: 12px; color: #56d364; }
.code-block { background: #0d1117; border: 1px solid #30363d; border-radius: 6px; padding: 12px; font-family: monospace; font-size: 12px; color: #c9d1d9; white-space: pre-wrap; overflow-x: auto; }
.method { display: inline-block; padding: 2px 6px; border-radius: 4px; font-size: 10px; font-weight: 600; margin-right: 6px; }
.method-post { background: #1a2a3a; color: #79c0ff; border: 1px solid #1f6feb; }
.method-get  { background: #1a3a1a; color: #56d364; border: 1px solid #238636; }
.path { font-family: monospace; color: #e6edf3; font-size: 13px; }
.desc { font-size: 12px; color: #8b949e; margin-top: 2px; }
.loading { color: #484f58; font-style: italic; font-size: 13px; }
.ep-row { padding: 8px 0; border-bottom: 1px solid #21262d; }
.ep-row:last-child { border-bottom: none; }

/* Collector */
.collector-controls { display: flex; gap: 10px; align-items: center; margin-bottom: 14px; flex-wrap: wrap; }
.collector-controls label { font-size: 13px; color: #8b949e; }
.count-input { background: #0d1117; border: 1px solid #30363d; border-radius: 6px; color: #e6edf3; font-size: 13px; padding: 5px 10px; width: 90px; outline: none; }
.count-input:focus { border-color: #1f6feb; }
.btn { padding: 6px 16px; border-radius: 6px; font-size: 13px; font-weight: 600; cursor: pointer; border: none; transition: opacity .15s; }
.btn:hover { opacity: .8; }
.btn:disabled { opacity: .4; cursor: not-allowed; }
.btn-green  { background: #238636; color: #fff; }
.btn-red    { background: #da3633; color: #fff; }
.btn-gray   { background: #21262d; color: #c9d1d9; border: 1px solid #30363d; }
.log-box { background: #0d1117; border: 1px solid #30363d; border-radius: 6px; padding: 10px 12px; font-family: monospace; font-size: 12px; color: #8b949e; height: 220px; overflow-y: auto; white-space: pre-wrap; word-break: break-all; }
.log-box .log-line { padding: 1px 0; line-height: 1.5; }
.log-box .log-ok   { color: #56d364; }
.log-box .log-err  { color: #f85149; }
.log-box .log-info { color: #79c0ff; }
.progress-bar-wrap { height: 4px; background: #21262d; border-radius: 2px; margin-bottom: 12px; overflow: hidden; }
.progress-bar { height: 100%; background: #1f6feb; border-radius: 2px; width: 0%; transition: width .3s; }
.progress-bar.indeterminate { animation: indeterminate 1.5s infinite; }
@keyframes indeterminate { 0%{width:0%;margin-left:0} 50%{width:60%;margin-left:20%} 100%{width:0%;margin-left:100%} }
.token-count-big { font-size: 28px; font-weight: 700; color: #56d364; }
.token-count-label { font-size: 12px; color: #8b949e; margin-top: 2px; }
</style>
</head>
<body>

<div class="header">
  <div class="dot" id="statusDot"></div>
  <div>
    <h1>GLM-ZAI-2API</h1>
    <div class="sub" id="modeSub">OpenAI-compatible proxy · loading…</div>
  </div>
</div>

<div class="grid">

  <!-- Status -->
  <div class="card">
    <div class="card-title">Status</div>
    <div class="row"><span class="row-label">Proxy</span><span class="badge badge-green" id="proxyBadge">running</span></div>
    <div class="row"><span class="row-label">Port</span><span class="row-value" id="portVal">—</span></div>
    <div class="row"><span class="row-label">Endpoint</span><span class="endpoint" id="endpointVal">—</span></div>
    <div class="row">
      <span class="row-label">Tokens in DB</span>
      <span class="row-value" id="tokenCountVal"><span class="loading">loading…</span></span>
    </div>
    <div class="row"><span class="row-label">User</span><span class="row-value" id="userVal">—</span></div>
    <div class="row"><span class="row-label">FE Version</span><span class="row-value" id="feVerVal" style="font-size:12px">—</span></div>
    <div class="row"><span class="row-label">Sessions</span><span class="row-value" id="sessionsVal">0</span></div>
  </div>

  <!-- Features -->
  <div class="card">
    <div class="card-title">Features <span style="font-size:10px;font-weight:400;color:#484f58;text-transform:none">(click to toggle)</span></div>
    <div class="row"><span class="row-label">Web Search</span><span class="badge badge-gray" id="featSearch" onclick="toggleFeat('webSearch',this)">OFF</span></div>
    <div class="row"><span class="row-label">Thinking</span><span class="badge badge-gray" id="featThink" onclick="toggleFeat('thinking',this)">OFF</span></div>
    <div class="row"><span class="row-label">Image Gen</span><span class="badge badge-gray" id="featImage" onclick="toggleFeat('imageGen',this)">OFF</span></div>
    <div class="row"><span class="row-label">Preview Mode</span><span class="badge badge-gray" id="featPreview" onclick="toggleFeat('previewMode',this)">OFF</span></div>
    <div class="row"><span class="row-label">Persist History</span><span class="badge badge-gray" id="featPersist" onclick="toggleFeat('persistHistory',this)">OFF</span></div>
  </div>

  <!-- Token Collector -->
  <div class="card" style="grid-column: 1 / -1;">
    <div class="card-title">Token Collector
      <span id="collectorBadge" class="badge badge-gray" style="margin-left:8px;text-transform:none;font-size:10px">idle</span>
    </div>

    <div class="collector-controls">
      <label for="countInput">Tokens to collect:</label>
      <input id="countInput" class="count-input" type="number" value="750" min="1" max="1250" />
      <button id="btnStart" class="btn btn-green" onclick="startCollect()">▶ Start Collection</button>
      <button id="btnStop"  class="btn btn-red"   onclick="stopCollect()" disabled>■ Stop</button>
      <span style="font-size:12px;color:#484f58">Max: 1250 · Uses headless Chromium inside the container</span>
    </div>

    <div class="progress-bar-wrap">
      <div class="progress-bar" id="progressBar"></div>
    </div>

    <div class="log-box" id="logBox">
      <div class="log-line" style="color:#484f58">Collector output will appear here…</div>
    </div>
  </div>

  <!-- API Endpoints -->
  <div class="card" style="grid-column: 1 / -1;">
    <div class="card-title">API Endpoints</div>
    <div class="ep-row">
      <span class="method method-post">POST</span>
      <span class="path">/v1/chat/completions</span>
      <div class="desc">OpenAI-compatible chat (stream / non-stream, tool calling)</div>
    </div>
    <div class="ep-row">
      <span class="method method-get">GET</span>
      <span class="path">/v1/models</span>
      <div class="desc">Available model list</div>
    </div>
    <div class="ep-row">
      <span class="method method-post">POST</span>
      <span class="path">/features</span>
      <div class="desc">Toggle webSearch, thinking, imageGen, previewMode, persistHistory</div>
    </div>
    <div class="ep-row">
      <span class="method method-post">POST</span>
      <span class="path">/admin/collect-tokens</span>
      <div class="desc">Start token collection (body: {"count": 750})</div>
    </div>
    <div class="ep-row">
      <span class="method method-get">GET</span>
      <span class="path">/admin/collect-tokens/stream</span>
      <div class="desc">SSE stream of collector progress logs</div>
    </div>
    <div class="ep-row">
      <span class="method method-post">POST</span>
      <span class="path">/admin/session/clear</span>
      <div class="desc">Clear all session history</div>
    </div>
    <div class="ep-row">
      <span class="method method-get">GET</span>
      <span class="path">/status</span>
      <div class="desc">Z.AI session status (public)</div>
    </div>
  </div>

  <!-- Test -->
  <div class="card" style="grid-column: 1 / -1;">
    <div class="card-title">Test the OpenAI Endpoint</div>
    <div class="code-block"># Non-streaming
curl -X POST http://__HOST__/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer __AUTH_TOKEN__" \
  -d '{"model":"glm-4.7","messages":[{"role":"user","content":"Hello!"}],"stream":false}'

# Streaming
curl -N -X POST http://__HOST__/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer __AUTH_TOKEN__" \
  -d '{"model":"glm-4.7","stream":true,"messages":[{"role":"user","content":"Say hi"}]}'</div>
  </div>

</div>

<script>
var AUTH = '__AUTH_TOKEN__';
var sseSource = null;
var collecting = false;

function setFeature(id, on) {
  var el = document.getElementById(id);
  if (!el) return;
  el.className = 'badge ' + (on ? 'badge-green' : 'badge-gray');
  el.textContent = on ? 'ON' : 'OFF';
}

function setCollecting(running) {
  collecting = running;
  var badge = document.getElementById('collectorBadge');
  var bar   = document.getElementById('progressBar');
  var btnS  = document.getElementById('btnStart');
  var btnSt = document.getElementById('btnStop');

  if (running) {
    badge.className = 'badge badge-yellow';
    badge.textContent = 'running';
    bar.className = 'progress-bar indeterminate';
    btnS.disabled = true;
    btnSt.disabled = false;
  } else {
    badge.className = 'badge badge-gray';
    badge.textContent = 'idle';
    bar.className = 'progress-bar';
    bar.style.width = '0%';
    btnS.disabled = false;
    btnSt.disabled = true;
  }
}

function appendLog(msg, cls) {
  var box = document.getElementById('logBox');
  // Clear placeholder on first real log
  if (box.children.length === 1 && box.children[0].style.color === 'rgb(72, 79, 88)') {
    box.innerHTML = '';
  }
  var line = document.createElement('div');
  line.className = 'log-line ' + (cls || '');
  var ts = new Date().toLocaleTimeString();
  line.textContent = '[' + ts + '] ' + msg;
  box.appendChild(line);
  box.scrollTop = box.scrollHeight;
}

function startCollect() {
  var count = parseInt(document.getElementById('countInput').value) || 750;
  if (count < 1) count = 1;
  if (count > 1250) count = 1250;

  document.getElementById('logBox').innerHTML = '';
  appendLog('Starting collection of ' + count + ' tokens…', 'log-info');
  setCollecting(true);

  fetch('/admin/collect-tokens', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + AUTH
    },
    body: JSON.stringify({count: count})
  }).then(function(r) { return r.json(); }).then(function(d) {
    if (d.error) {
      appendLog('Error: ' + d.error, 'log-err');
      setCollecting(false);
      return;
    }
    appendLog('Collection started → ' + d.dbPath, 'log-info');
    startSSE();
  }).catch(function(e) {
    appendLog('Request failed: ' + e, 'log-err');
    setCollecting(false);
  });
}

function stopCollect() {
  fetch('/admin/collect-tokens/stop', {
    method: 'POST',
    headers: {'Authorization': 'Bearer ' + AUTH}
  });
  appendLog('Stop requested…', 'log-info');
  if (sseSource) { sseSource.close(); sseSource = null; }
  setCollecting(false);
}

function startSSE() {
  if (sseSource) { sseSource.close(); }
  sseSource = new EventSource('/admin/collect-tokens/stream?auth=' + encodeURIComponent(AUTH));

  sseSource.onmessage = function(e) {
    try {
      var d = JSON.parse(e.data);
      if (d.done) {
        appendLog('✓ Done!', 'log-ok');
        setCollecting(false);
        sseSource.close();
        sseSource = null;
        refresh();
        return;
      }
      if (d.log) {
        var cls = '';
        if (d.log.startsWith('✓') || d.log.indexOf('Saved') >= 0 || d.log.indexOf('complete') >= 0) cls = 'log-ok';
        else if (d.log.startsWith('ERROR') || d.log.startsWith('✗')) cls = 'log-err';
        else if (d.log.startsWith('Starting') || d.log.startsWith('Attempt') || d.log.startsWith('Navigating')) cls = 'log-info';
        appendLog(d.log, cls);
      }
      if (typeof d.tokenCount !== 'undefined') {
        document.getElementById('tokenCountVal').textContent = d.tokenCount + ' tokens';
      }
      if (typeof d.running !== 'undefined' && !d.running && collecting) {
        setCollecting(false);
      }
    } catch(err) {}
  };

  sseSource.onerror = function() {
    if (collecting) {
      appendLog('SSE connection lost, retrying…', 'log-err');
    }
  };
}

// SSE auth via query param — server reads it
// Override handleCollectTokensStream to accept ?auth= too
// (already handled server-side via withAuth which checks header OR query)

async function refresh() {
  try {
    var res = await fetch('/status');
    var d = await res.json();

    document.getElementById('modeSub').textContent = 'OpenAI-compatible proxy · ' + (d.mode || 'direct');
    var dot = document.getElementById('statusDot');
    var proxyBadge = document.getElementById('proxyBadge');
    if (d.connected) {
      dot.style.background = '#56d364';
      dot.style.boxShadow = '0 0 8px #56d364';
      proxyBadge.className = 'badge badge-green';
      proxyBadge.textContent = 'running';
    } else {
      dot.style.background = '#f85149';
      dot.style.boxShadow = '0 0 8px #f85149';
      proxyBadge.className = 'badge badge-red';
      proxyBadge.textContent = 'offline';
    }

    document.getElementById('portVal').textContent = d.port || '—';
    document.getElementById('endpointVal').textContent = 'http://127.0.0.1:' + (d.port || '5082') + '/v1';
    document.getElementById('tokenCountVal').textContent = (d.tokenCount != null ? d.tokenCount + ' tokens' : '—');
    document.getElementById('userVal').textContent = d.userName || '—';
    document.getElementById('feVerVal').textContent = d.feVersion || '—';
    document.getElementById('sessionsVal').textContent = d.activeSessions || 0;

    var f = d.features || {};
    setFeature('featSearch', f.webSearch);
    setFeature('featThink', f.thinking);
    setFeature('featImage', f.imageGen);
    setFeature('featPreview', f.previewMode);
    setFeature('featPersist', f.persistHistory);

    if (d.collectorRunning && !collecting) {
      setCollecting(true);
      startSSE();
    }
  } catch(e) { console.error('dashboard refresh error:', e); }
}

function toggleFeat(key, el) {
  var isOn = el.textContent === 'ON';
  var body = {}; body[key] = !isOn;
  fetch('/features', {
    method: 'POST',
    headers: {'Content-Type':'application/json','Authorization':'Bearer ' + AUTH},
    body: JSON.stringify(body)
  }).then(function(r){return r.json()}).then(function(d){
    if(d.success) refresh();
  }).catch(function(e){console.error('toggle error:',e)});
}

refresh();
setInterval(refresh, 5000);
</script>
</body>
</html>`