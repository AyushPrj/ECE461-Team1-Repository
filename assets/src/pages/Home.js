// src/pages/Home.js

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
  const handleResetRegistry = () => {
    console.log('Reset registry');
  };

  return (
    <Container>
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
        <Grid item xs={12} sm={6} md={4}>
          <Button
            component={RouterLink}
            to="/create"
            variant="contained"
            fullWidth
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
          >
            Reset Registry
          </Button>
        </Grid>
      </Grid>
    </Container>
  );
};

export default Home;
