import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080';

export default {
  authenticate: (token) => axios.put(`${API_BASE_URL}/authenticate`, token),
  createPackage: (packageData) => axios.post(`${API_BASE_URL}/package`, packageData),
  getPackageById: (id) => axios.get(`${API_BASE_URL}/package/${id}`),
  updatePackage: (id, packageData) => axios.put(`${API_BASE_URL}/package/${id}`, packageData),
  deletePackage: (id) => axios.delete(`${API_BASE_URL}/package/${id}`),
  ratePackage: (id) => axios.get(`${API_BASE_URL}/package/${id}/rate`),
  listPackages: (queryParams) => axios.post(`${API_BASE_URL}/packages`, queryParams),
  resetRegistry: () => axios.delete(`${API_BASE_URL}/reset`),
  getPackageByName: (name) => axios.get(`${API_BASE_URL}/package/byName/${name}`),
  deletePackageByName: (name) => axios.delete(`${API_BASE_URL}/package/byName/${name}`),
  getPackageByRegEx: (regex) => axios.post(`${API_BASE_URL}/package/byRegEx`, regex),
};
