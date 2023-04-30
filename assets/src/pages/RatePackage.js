import React, { useState } from 'react';
import {
    Container,
    Typography,
    Grid,
    Box,
    FormControl,
    InputLabel,
    OutlinedInput,
    Button,
    Paper,
} from '@mui/material';
import { useForm } from 'react-hook-form';

const RatePackage = () => {
    const { register, handleSubmit, formState: { errors } } = useForm();
    const [response, setResponse] = useState(null);

    const onSubmit = async (formData) => {
        console.log("Form errors:", errors);
        console.log("Form data:", formData);

        const { id } = formData;

        try {
            const response = await fetch(`https://webservice-381819.uc.r.appspot.com/package/${id}/rate`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': '*/*',
                    'X-Authorization': 'hello23'
                },
            });

            if (response.ok) {
                const responseData = await response.json();
                console.log(responseData);
                setResponse(responseData);
            } else {
                const errorMessage = await response.text();
                console.error('Error:', errorMessage);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    };

    return (
        <Container component="main" tabIndex="0">
            <Box mt={4} mb={4}>
                <Typography variant="h4" component="h1" gutterBottom>
                    Rate Package
                </Typography>
            </Box>
            <form onSubmit={handleSubmit(onSubmit)}>
                <Grid container spacing={2}>
                    <Grid item xs={12}>
                        <FormControl fullWidth variant="outlined">
                            <InputLabel htmlFor="id">Package ID</InputLabel>
                            <OutlinedInput
                                id="id"
                                type="text"
                                label="Package ID"
                                {...register('id', { required: true })}
                                required
                                autoFocus
                                tabIndex="0"
                                inputProps={{
                                    'aria-label': 'Package ID',
                                }}
                            />
                        </FormControl>
                    </Grid>
                    <Grid item xs={12}>
                        <Button
                            variant="contained"
                            color="primary"
                            type="submit"
                            tabIndex="0"
                            aria-label="Submit package ID"
                        >
                            Submit
                        </Button>
                    </Grid>
                </Grid>
            </form>
            {response && (
                <Box mt={4}>
                    <Paper elevation={3} sx={{ padding: 2 }}>
                        <pre tabIndex="0" aria-label="Package rating response">
                            {JSON.stringify(response, null, 2)}
                        </pre>
                    </Paper>
                </Box>
            )}
        </Container>
    );
};

export default RatePackage;
