// dashboard.go — embedded HTML dashboard (ports the getDashboardHTML template).
package main

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Z.AI Direct Bridge</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      background: linear-gradient(135deg, #1e3a5f 0%, #0d1b2a 50%, #1b263b 100%);
      min-height: 100vh; color: #e0e0e0; padding: 20px;
    }
    .container { max-width: 1200px; margin: 0 auto; }
    .header {
      text-align: center; padding: 40px 20px;
      background: rgba(255,255,255,0.05); border-radius: 16px;
      margin-bottom: 30px; border: 1px solid rgba(255,255,255,0.1);
    }
    .header h1 {
      font-size: 2.5rem;
      background: linear-gradient(135deg, #3b82f6, #1d4ed8, #60a5fa);
      -webkit-background-clip: text; -webkit-text-fill-color: transparent;
      margin-bottom: 10px;
    }
    .header p { color: #888; font-size: 1.1rem; }
    .badges { display: flex; gap: 8px; justify-content: center; margin-top: 12px; flex-wrap: wrap; }
    .badge {
      display: inline-block; padding: 4px 12px; border-radius: 12px;
      font-size: 0.8rem; font-weight: 700;
    }
    .badge-green { background: #22c55e; color: #000; }
    .badge-blue  { background: #3b82f6; color: #fff; }
    .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
    .card {
      background: rgba(255,255,255,0.05); border-radius: 12px;
      padding: 24px; border: 1px solid rgba(255,255,255,0.1);
    }
    .card h2 { color: #60a5fa; margin-bottom: 16px; font-size: 1.2rem; }
    .stat-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 12px; }
    .stat { background: rgba(0,0,0,0.2); padding: 12px; border-radius: 8px; }
    .stat .label { color: #888; font-size: 0.85rem; }
    .stat .value { color: #60a5fa; font-weight: 600; font-size: 1.5rem; margin-top: 4px; }
    .code-block {
      background: #0d1117; border-radius: 8px; padding: 16px; overflow-x: auto;
      font-family: 'Monaco', 'Menlo', monospace; font-size: 0.85rem;
      border: 1px solid #30363d; margin: 12px 0;
    }
    .code-block code { color: #c9d1d9; white-space: pre-wrap; }
    .endpoint { background: rgba(0,0,0,0.2); padding: 12px; border-radius: 8px; margin-bottom: 8px; }
    .method {
      display: inline-block; padding: 4px 8px; border-radius: 4px;
      font-size: 0.75rem; font-weight: 600; margin-right: 8px;
    }
    .method.get { background: #22c55e; color: #000; }
    .method.post { background: #3b82f6; color: #fff; }
    .path { font-family: monospace; color: #e0e0e0; }
    .desc { color: #888; font-size: 0.85rem; margin-top: 4px; }
    .section-label {
      font-size: 0.75rem; font-weight: 700; text-transform: uppercase;
      letter-spacing: 0.1em; color: #a855f7; margin: 16px 0 8px;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Z.AI Direct Bridge</h1>
      <p>HTTP-only mode — No browser required (Go)</p>
      <div class="badges">
        <span class="badge badge-green">Direct Mode</span>
        <span class="badge badge-blue">OpenAI Compatible</span>
      </div>
    </div>

    <div class="grid">
      <div class="card">
        <h2>Session Status</h2>
        <div class="stat-grid">
          <div class="stat">
            <div class="label">Connection</div>
            <div class="value" id="sessionStatus">...</div>
          </div>
          <div class="stat">
            <div class="label">User</div>
            <div class="value" id="sessionUser" style="font-size:1rem">...</div>
          </div>
          <div class="stat">
            <div class="label">Messages</div>
            <div class="value" id="msgCount">0</div>
          </div>
          <div class="stat">
            <div class="label">FE Version</div>
            <div class="value" id="feVersion" style="font-size:0.85rem">...</div>
          </div>
        </div>
      </div>

      <div class="card">
        <h2>Features</h2>
        <div class="stat-grid">
          <div class="stat"><div class="label">Web Search</div><div class="value" id="featSearch">-</div></div>
          <div class="stat"><div class="label">Thinking</div><div class="value" id="featThink">-</div></div>
          <div class="stat"><div class="label">Image Gen</div><div class="value" id="featImage">-</div></div>
          <div class="stat"><div class="label">Preview</div><div class="value" id="featPreview">-</div></div>
        </div>
      </div>

      <div class="card" style="grid-column: span 2;">
        <h2>API Endpoints</h2>

        <div class="section-label">OpenAI-Compatible</div>
        <div class="endpoint">
          <span class="method post">POST</span>
          <span class="path">/v1/chat/completions</span>
          <div class="desc">OpenAI-compatible chat endpoint. Supports streaming.</div>
        </div>
        <div class="endpoint">
          <span class="method get">GET</span>
          <span class="path">/v1/models</span>
          <div class="desc">Model list</div>
        </div>

        <div class="section-label">Management</div>
        <div class="endpoint">
          <span class="method post">POST</span>
          <span class="path">/features</span>
          <div class="desc">Toggle webSearch, thinking, imageGen, previewMode, persistHistory</div>
        </div>
        <div class="endpoint">
          <span class="method post">POST</span>
          <span class="path">/admin/session/clear</span>
          <div class="desc">Clear conversation history</div>
        </div>
      </div>

      <div class="card" style="grid-column: span 2;">
        <h2>Test the OpenAI endpoint</h2>
        <div class="code-block">
          <code># Non-streaming
curl -X POST http://__HOST__/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer __AUTH_TOKEN__" \
  -d '{"model":"glm-4.7","messages":[{"role":"user","content":"Hello!"}],"stream":false}'

# Streaming
curl -X POST http://__HOST__/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer __AUTH_TOKEN__" \
  -d '{"model":"glm-4.7","stream":true,"messages":[{"role":"user","content":"Say hi"}]}'</code>
        </div>
      </div>
    </div>
  </div>

  <script>
    async function updateStatus() {
      try {
        const res = await fetch('/status');
        const d = await res.json();
        document.getElementById('sessionStatus').textContent = d.connected ? 'OK' : 'Off';
        document.getElementById('sessionUser').textContent = d.userName || '-';
        document.getElementById('msgCount').textContent = d.activeSessions;
        document.getElementById('feVersion').textContent = d.feVersion || '-';
        document.getElementById('featSearch').textContent = d.features?.webSearch ? 'ON' : 'OFF';
        document.getElementById('featThink').textContent = d.features?.thinking ? 'ON' : 'OFF';
        document.getElementById('featImage').textContent = d.features?.imageGen ? 'ON' : 'OFF';
        document.getElementById('featPreview').textContent = d.features?.previewMode ? 'ON' : 'OFF';
      } catch(e) { console.error(e); }
    }
    updateStatus();
    setInterval(updateStatus, 3000);
  </script>
</body>
</html>`
