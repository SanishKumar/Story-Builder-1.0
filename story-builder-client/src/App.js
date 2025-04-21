import React, { useState, useEffect } from 'react';

function App() {
  const [prompt, setPrompt] = useState("");
  const [generatedText, setGeneratedText] = useState("");
  const [story, setStory] = useState([]);

  const handleGenerate = async () => {
    const response = await fetch("http://localhost:8080/generate", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ prompt }),
    });
    const data = await response.json();
    setGeneratedText(data.generated_text);
  };

  // You can also implement real-time updates using SSE or WebSockets to listen for story updates.
  // For example, connect via EventSource to a /stream endpoint.

  return (
    <div>
      <h1>Generative AI Story Builder</h1>
      <input
        type="text"
        value={prompt}
        onChange={(e) => setPrompt(e.target.value)}
        placeholder="Enter a story prompt"
      />
      <button onClick={handleGenerate}>Generate Segment</button>
      <p>Generated: {generatedText}</p>
      <h2>Current Story</h2>
      <ul>
        {story.map((seg) => (
          <li key={seg.id}>{seg.text} (Votes: {seg.votes})</li>
        ))}
      </ul>
    </div>
  );
}

export default App;
