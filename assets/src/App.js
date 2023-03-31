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
      const formData = new FormData();
      formData.append('userInput', inputValue);
      //formData.append('zipFile', zipFile);
      const zip = await JSZip.loadAsync(await zipFile);
      const files = zip.files;
      const fileData = {};
      var fileName;
      for (const fileName in files) {
        const file = files[fileName];
        const fileContent = await file.async('uint8array');
        fileData[fileName] = fileContent;
      }
      fileName = files[0]
      
      let zipfilename = zipFile.name;
      
      const nameArray = zipfilename.split("-");
      zipfilename = zipfilename.substring(0, zipfilename.length - nameArray[nameArray.length - 1].length - 1)

      let name = "cloudinary/cloudinary-video-player"

      // const res0 = await fetch('http://localhost:5500/raterepo', {
      //   method: 'GET',
      //   body: urlname,
      // });

      const res = await fetch('http://localhost:5500/repo', {
        method: 'POST',
        //body: "{\r\n    \"name\": \"cloudinar5\",\r\n    \"rampup\": 0.23,\r\n    \"correctness\": 1,\r\n    \"responsivemaintainer\": 0.5,\r\n    \"busfactor\": 0.4,\r\n    \"reviewcoverage\": 0.2,\r\n    \"dependancypinning\": 0.6,\r\n    \"license\": 1,\r\n    \"net\": 0.8\r\n}",
        body: `{"url": "${name}"}`
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
