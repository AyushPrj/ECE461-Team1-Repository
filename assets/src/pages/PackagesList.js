import React, { useState } from 'react';
import {
  Container,
  Typography,
  Box,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Grid,
  Button,
} from '@mui/material';

const PackagesList = () => {
  const [packages, setPackages] = useState([]);
  const [pageNumber, setPageNumber] = useState('');
  const [name, setName] = useState('');
  const [version, setVersion] = useState('');

  const fetchPackages = async () => {
    try {
      const response = await fetch('https://webservice-381819.uc.r.appspot.com/packages', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': '*/*',
        },
        body: JSON.stringify([
          {
            Version: version,
            Name: name,
          },
        ]),
      });

      if (response.ok) {
        const responseData = await response.json();
        console.log(responseData);
        setPackages(responseData);
      } else {
        const errorMessage = await response.text();
        console.error('Error:', errorMessage);
      }
    } catch (error) {
      console.error('Error:', error);
    }
  };

  const handleFetchPackages = () => {
    fetchPackages();
  };

  return (
    <Container component="main" tabIndex="0">
      <Box mt={4} mb={4}>
        <Typography variant="h4" component="h1" gutterBottom>
          Packages List
        </Typography>
      </Box>
      <Grid container spacing={2}>
        {/* Add tabIndex="0" and aria-label to input fields */}
        <Grid item xs={12} sm={4}>
          <TextField
            label="Page Number"
            value={pageNumber}
            onChange={(e) => setPageNumber(e.target.value)}
            fullWidth
            tabIndex="0"
            aria-label="Page Number"
          />
        </Grid>
        <Grid item xs={12} sm={4}>
          <TextField
            label="Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            fullWidth
            tabIndex="0"
            aria-label="Package Name"
          />
        </Grid>
        <Grid item xs={12} sm={4}>
          <TextField
            label="Version"
            value={version}
            onChange={(e) => setVersion(e.target.value)}
            fullWidth
            tabIndex="0"
            aria-label="Package Version"
          />
        </Grid>
        <Grid item xs={12}>
          <Button
            variant="contained"
            color="primary"
            onClick={handleFetchPackages}
            tabIndex="0"
            aria-label="Fetch Packages"
          >
            Fetch Packages
          </Button>
        </Grid>
      </Grid>
      <TableContainer component={Paper}>
        <Table aria-label="Packages List Table">
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>Version</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {packages.map((pkg) => (
              <TableRow key={pkg.ID}>
                <TableCell component="th" scope="row">
                  {pkg.ID}
                </TableCell>
                <TableCell>{pkg.Name}</TableCell>
                <TableCell>{pkg.Version}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
};

export default PackagesList;
