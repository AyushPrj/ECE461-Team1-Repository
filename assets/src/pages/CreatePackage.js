import React from 'react';
import {
    Container,
    Typography,
    Grid,
    Box,
    FormControl,
    InputLabel,
    OutlinedInput,
    InputAdornment,
    IconButton,
    Button
} from '@mui/material';
import { useForm } from 'react-hook-form';
import { AttachFile } from '@mui/icons-material';
import JSZip from 'jszip';

const CreatePackage = () => {
    const { register, handleSubmit, setValue, getValues } = useForm();

    const handleUpload = async (event) => {
        const file = event.target.files[0];
        const reader = new FileReader();

        reader.onloadend = async () => {
            const base64 = reader.result.split(',')[1]; // Remove the prefix
            console.log(base64);
            setValue('Content', base64);

            const zip = new JSZip();
            const loadedZip = await zip.loadAsync(file);

            const folderName = Object.keys(loadedZip.files)[0];
            const packageJsonPath = `${folderName}package.json`;
            const packageJsonFile = await loadedZip.file(packageJsonPath).async('string');

            const packageJson = JSON.parse(packageJsonFile);

            if (packageJson.repository && packageJson.repository.url) {
                console.log("HELLO3")

                let url = packageJson.repository.url;

                if (url.endsWith('.git')) {
                    url = url.slice(0, -4);
                }
                setValue('url', url);
                console.log(url);

            }

            console.log("HELLO4")

            setValue('jsprogram', 'arbitraryValue');

        };


        if (file) {
            reader.readAsDataURL(file);
        }
    };

    const onSubmit = async (formData) => {
        const { Content, url, jsprogram } = formData;

        const requestData = {
            Content: Content,
            URL: url,
            JSProgram: jsprogram,
        };
        console.log(requestData);

        try {
            const response = await fetch('https://webservice-381819.uc.r.appspot.com/package', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': '*/*',
                    // If you need an authentication token, add it here
                    // 'Authorization': 'Bearer ' + token
                },
                body: JSON.stringify(requestData),
            });

            if (response.ok) {
                const responseData = await response.json();
                console.log(responseData);
                // Handle successful API call, e.g., show a success message or redirect to another page
            } else {
                const errorMessage = await response.text();
                console.error('Error:', errorMessage);
                // Handle error, e.g., show an error message to the user
            }
        } catch (error) {
            console.error('Error:', error);
            // Handle network error, e.g., show an error message to the user
        }
    };


    return (
        <Container>
            <Box mt={4} mb={4}>
                <Typography variant="h4" component="h1" gutterBottom>
                    Create Package
                </Typography>
            </Box>
            <form onSubmit={handleSubmit(onSubmit)}>
                <Grid container spacing={2}>
                    <Grid item xs={12}>
                        <FormControl fullWidth variant="outlined">
                            <InputLabel htmlFor="content">Content (ZIP File)</InputLabel>
                            <OutlinedInput
                                id="content"
                                type="file"
                                accept=".zip"
                                onChange={handleUpload}
                                label="Content (ZIP File)"
                                inputProps={{
                                    ref: register('Content', { required: true })
                                }}
                                startAdornment={
                                    <InputAdornment position="start">
                                        <IconButton edge="start">
                                            <AttachFile />
                                        </IconButton>
                                    </InputAdornment>
                                }
                            />
                        </FormControl>
                    </Grid>
                    <Grid item xs={12}>
                        <Button variant="contained" color="primary" type="submit">
                            Submit
                        </Button>
                    </Grid>
                </Grid>
            </form>
        </Container>
    );
};

export default CreatePackage;
