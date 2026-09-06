import './style.css';
import { SelectFolder,AnalyzeProject } from '../wailsjs/go/main/App';
document.documentElement.lang = 'en';
document.title = 'Cleanupcheck';




// Render the interface before attaching button handlers.
document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
  <div class="workspace">
    <header class="topbar">
      <div class="brand"><span class="brand-icon" aria-hidden="true">✓</span>
        <span>cleanupcheck<small>GO CODE ANALYSIS</small></span></div>
      <span class="version">Desktop · Preview</span>
    </header>
    <main>
      <section class="intro" aria-labelledby="page-title">
        <span class="eyebrow">CHECK YOUR CODE</span>
        <h1 id="page-title">Find unclosed resources.</h1>
        <p>Check your Go project for unclosed HTTP response bodies.
          Review findings with file names and line numbers in one place.</p>
      </section>
      <section class="panel project-panel" aria-labelledby="project-title">
        <div class="section-heading">
          <div><h2 id="project-title">Select a project</h2>
            <p>Choose the Go project folder you want to analyze.</p></div>
          <span class="step" aria-hidden="true">01</span>
        </div>
        <label for="project-path">Project folder</label>
        <div class="folder-row">
          <input id="project-path" type="text" placeholder="No folder selected"
            readonly aria-describedby="connection-note" />
          <button id="select-folder" class="button secondary"
            aria-describedby="connection-note">Choose Folder</button>
        </div>
        <div class="project-footer">
          <p id="connection-note" role="status">Choose a Go project folder. Analysis is not available yet.</p>
          <button id="analyze" class="button primary" disabled
            aria-describedby="connection-note">Analyze <span aria-hidden="true">→</span></button>
        </div>
      </section>
      <section class="panel results-panel" aria-labelledby="results-title">
        <div class="section-heading results-heading">
          <h2 id="results-title">Analysis results</h2>
          <span class="status">Not run yet</span>
        </div>
        <div id="results" aria-live="polite">
          <div class="empty-state">
            <span class="empty-icon" aria-hidden="true">⌕</span>
            <h3>Waiting for your first analysis</h3>
            <p>Findings will appear here once the analysis is complete.</p>
          </div>
        </div>
      </section>
      <aside class="tip"><p><strong>What do we check?</strong> Whether HTTP response bodies are closed
        with <code>resp.Body.Close()</code> after a request.</p></aside>
    </main>
    <footer>Cleanupcheck <span>Resource cleanup checks for Go projects</span></footer>
  </div>
`;

const folderButton = document.querySelector<HTMLButtonElement>('#select-folder')!;
const projectPath = document.querySelector<HTMLInputElement>('#project-path')!;
const connectionNote = document.querySelector<HTMLParagraphElement>('#connection-note')!;

// Set up handlers only after the interface exists.
const analyzeButton =
  document.querySelector<HTMLButtonElement>('#analyze')!;

const results =
  document.querySelector<HTMLDivElement>('#results')!;

analyzeButton.addEventListener('click', async () => {
  if (!projectPath.value) return;

  analyzeButton.disabled = true;
  folderButton.disabled = true;
  analyzeButton.textContent = 'Analyzing…';
  results.textContent = 'Checking your project…';

  try {
    const findings = await AnalyzeProject(projectPath.value);

    results.replaceChildren();

    if (findings.length === 0) {
      results.textContent = 'No issues found.';
    }

    for (const finding of findings) {
      const row = document.createElement('p');

      row.textContent =
        `${finding.file}:${finding.line} — ${finding.message}`;

      results.appendChild(row);
    }
  } catch (error) {
    results.textContent = `Analysis failed: ${String(error)}`;
  } finally {
    analyzeButton.disabled = false;
    folderButton.disabled = false;
    analyzeButton.textContent = 'Analyze';
  }
});

folderButton.addEventListener('click', async () => {
  // Prevent opening multiple dialogs while waiting for the user's selection.
  folderButton.disabled = true;
  folderButton.textContent = 'Choosing…';

  try {
    const folder = await SelectFolder();

    // Cancelling the dialog preserves any previously selected folder.
    if (folder) {
      projectPath.value = folder;
      projectPath.title = folder;
      analyzeButton.disabled = false;
    }
    connectionNote.textContent = projectPath.value
      ? 'Folder selected. Analysis is not available yet.'
      : 'No folder selected. Choose a Go project folder.';
  } catch (error) {
    connectionNote.textContent = `Could not open folder selection: ${String(error)}`;
  } finally {
    folderButton.disabled = false;
    folderButton.textContent = 'Choose Folder';
  }
});

