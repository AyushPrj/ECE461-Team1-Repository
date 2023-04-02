const JSZip = require('jszip');

async function readUrlFromZip(zipFileName) {
  try {
    console.log('Loading zip file...');
    const zipFile = new File([await fetch(zipFileName).then(res => res.blob())], zipFileName);
    console.log('Loading zip archive...');
    const zip = await JSZip.loadAsync(zipFile);
    console.log('Extracting package.json...');
    const packageFile = await zip.file('package.json').async('string');
    const packageJson = JSON.parse(packageFile);
    const url = packageJson.url;
    console.log('URL extracted:', url);
    return url;
  } catch (error) {
    console.error(error);
    throw error;
  }
}




module.exports = {
  readUrlFromZip
};