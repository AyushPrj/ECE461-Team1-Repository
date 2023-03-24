import React, { useState } from 'react';

function App() {
  const [response, setResponse] = useState(null);
  const [inputValue, setInputValue] = useState('');
  const [repoId, setRepoId] = useState('');
  const [zipFile, setZipFile] = useState(null);

  const handleInputChange = (event) => {
    setInputValue(event.target.value);
  };

  const handleZipFileChange = (event) => {
    setZipFile(event.target.files[0]);
  };

  const handleCreateRepoClick = async () => {
    try {
      const formData = new FormData();
      formData.append('userInput', inputValue);
      formData.append('zipFile', zipFile);

      const res = await fetch('http://localhost:5500/repo', {
        method: 'POST',
        body: formData,
      });

      const json = await res.json();
      setResponse(json);
    } catch (error) {
      console.error(error);
    }
  };

  const handleGetRepoClick = async () => {
    try {
      const res = await fetch(`http://localhost:5500/repo/${repoId}`);
      const json = await res.json();
      setResponse(json);
    } catch (error) {
      console.error(error);
    }
  };

  const handleEditRepoClick = async () => {
    try {
      const res = await fetch(`http://localhost:5500/repo/${repoId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ userInput: inputValue }),
      });
      const json = await res.json();
      setResponse(json);
    } catch (error) {
      console.error(error);
    }
  };

  const handleDeleteRepoClick = async () => {
    try {
      const res = await fetch(`http://localhost:5500/repo/${repoId}`, {
        method: 'DELETE',
      });
      const json = await res.json();
      setResponse(json);
    } catch (error) {
      console.error(error);
    }
  };

  const handleGetAllReposClick = async () => {
    try {
      const res = await fetch('http://localhost:5500/repos');
      const json = await res.json();
      setResponse(json);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div>
      <h1>My Website</h1>
      <form>
        <label>
          User Input:
          <input type="text" value={inputValue} onChange={handleInputChange} />
        </label>
        <br />
        <label>
          Zip File:
          <input type="file" onChange={handleZipFileChange} />
        </label>
        <br />
        <label>
          Repo ID:
          <input type="text" value={repoId} onChange={(event) => setRepoId(event.target.value)} />
        </label>
        <br />
        <button type="button" onClick={handleCreateRepoClick}>
          Create Repo
        </button>
        <button type="button" onClick={handleGetRepoClick}>
          Get Repo
        </button>
        <button type="button" onClick={handleEditRepoClick}>
          Edit Repo
        </button>
        <button type="button" onClick={handleDeleteRepoClick}>
          Delete Repo
        </button>
        <button type="button" onClick={handleGetAllReposClick}>
          Get All Repos
        </button>
      </form>
      {response && <pre>{JSON.stringify(response, null, 2)}</pre>}
    </div>
  );
}

export default App;
