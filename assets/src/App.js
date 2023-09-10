import React, { useState } from 'react';
const JSZip = require('jszip');

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
      console.log("START");

      const reader = new FileReader();

      reader.addEventListener('load', async () => {
        const buffer = reader.result;

        const zip = new JSZip();
        await zip.loadAsync(buffer);

        // get the first folder in the zip archive
        const firstFolderName = Object.keys(zip.files).find(name => {
          return name.endsWith('/') && name.split('/').length === 2;
        });

        if (!firstFolderName) {
          throw new Error('No folder found in zip archive');
        }

        // access the package.json file within the first folder
        const packageJsonText = await zip.file(`${firstFolderName}package.json`).async('text');

        // parse the JSON string and access the url field within the repository field
        const packageJson = JSON.parse(packageJsonText);
        let url = packageJson.repository.url;

        // remove anything before and including the third /
        const urlParts = url.split('/');
        urlParts.splice(0, 3);
        url = urlParts.join('/');

        // strip the last 4 characters from the url
        url = url.slice(0, -4);

        console.log("Fetch");
        const res = await fetch('http://localhost:5500/repo', {
          method: 'POST',
          body: `{"url": "${url}"}`
        });

        const json = await res.json();
        setResponse(json);

        console.log(url);
      });

      reader.readAsArrayBuffer(zipFile);
    } catch (err) {
      console.error(err);
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
