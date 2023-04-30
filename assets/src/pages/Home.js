import React from 'react';
import { Link as RouterLink } from 'react-router-dom';
import {
  Button,
  Container,
  Typography,
  Box,
  Grid,
  Link,
} from '@mui/material';

const Home = () => {
  const handleResetRegistry = async () => {
    try {
      const response = await fetch('https://webservice-381819.uc.r.appspot.com/reset', {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'Accept': '*/*',
          // If you need an authentication token, add it here
          // 'Authorization': 'Bearer ' + token
        },
      });

      if (response.ok) {
        const responseData = await response.json();
        console.log(responseData);
        // Handle successful API call, e.g., show a success message or update the UI
      } else {
        const errorMessage = await response.text();
        console.error('Error:', errorMessage);
        // Handle error, e.g., show an error message to the user
      }
    } catch (error) {
      console.error('Error:', error);
      // Handle network error, e.g., show an error message to the user
    }
  }

  return (
    <Container component="main" tabIndex="0">
      <Box mt={4} mb={4}>
        <Typography variant="h4" component="h1" gutterBottom>
          Welcome to the Package Manager
        </Typography>
        <Typography variant="body1">
          This application allows you to search, create, update, and delete
          packages in the package registry.
        </Typography>
      </Box>
      <Grid container spacing={2}>
        {/* Each Button has tabIndex="0" and an aria-label */}
        <Grid item xs={12} sm={6} md={4}>
          <Button
            component={RouterLink}
            to="/create"
            variant="contained"
            fullWidth
            tabIndex="0"
            aria-label="Create a Package"
          >
            Create a Package
          </Button>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Button
            component={RouterLink}
            to="/search"
            variant="contained"
            fullWidth
            tabIndex="0"
            aria-label="Search Packages"
          >
            Search Packages
          </Button>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Button
            component={RouterLink}
            to="/packages"
            variant="contained"
            fullWidth
            tabIndex="0"
            aria-label="View All Packages"
          >
            View All Packages
          </Button>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Button
            component={RouterLink}
            to="/update/:id"
            variant="contained"
            fullWidth
            tabIndex="0"
            aria-label="Update a Package"
          >
            Update a Package
          </Button>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Button
            component={RouterLink}
            to="/package/:id"
            variant="contained"
            fullWidth
            tabIndex="0"
            aria-label="View Package Details"
          >
            View Package Details
          </Button>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Button
            component={RouterLink}
            to="/package/:id/rate"
            variant="contained"
            fullWidth
            tabIndex="0"
            aria-label="Rate a Package"
          >
            Rate a Package
          </Button>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Button
            variant="contained"
            color="error"
            fullWidth
            onClick={handleResetRegistry}
            tabIndex="0"
            aria-label="Reset Registry"
          >
            Reset Registry
          </Button>
        </Grid>
      </Grid>
    </Container>
  );
};

export default Home;