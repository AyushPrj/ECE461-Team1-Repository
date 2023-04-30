import React, { useState } from 'react';
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

const UpdatePackage = () => {
    const { register, handleSubmit, setValue, getValues, trigger, formState: { errors } } = useForm();
    const [fileLoading, setFileLoading] = useState(false);
    const [packageName, setPackageName] = useState('');
    const [packageVersion, setPackageVersion] = useState('');

    const handleUpload = async (event) => {
        console.log("uploaded")
        const file = event.target.files[0];
        const reader = new FileReader();

        setFileLoading(true); // Add this line
        reader.onloadend = async () => {
            const base64 = reader.result.split(',')[1]; // Remove the prefix
            setValue('Content', base64);

            const zip = new JSZip();
            const loadedZip = await zip.loadAsync(file);

            const folderName = Object.keys(loadedZip.files)[0];
            const packageJsonPath = `${folderName}package.json`;
            const packageJsonFile = await loadedZip.file(packageJsonPath).async('string');

            const packageJson = JSON.parse(packageJsonFile);

            // Set Name and Version from package.json
            setPackageName(packageJson.name);
            setPackageVersion(packageJson.version);

            if (packageJson.repository && packageJson.repository.url) {
                let url = packageJson.repository.url;

                if (url.endsWith('.git')) {
                    url = url.slice(0, -4);
                }
                setValue('url', url);
            }

            setValue('jsprogram', 'arbitraryUpdate3');
            setFileLoading(false); // Add this line
        };

        if (file) {
            reader.readAsDataURL(file);
        }
        console.log("uploaded2")
    };


    const onSubmit = async (formData) => {
        console.log("Form errors:", errors);
        console.log("Form data:", formData);
        console.log("Hello");

        const { id, Content, url, jsprogram } = formData;

        console.log(packageName);
        console.log(packageVersion);

        const requestData = {
            metadata: {
                Name: packageName,
                Version: packageVersion,
                ID: id,
            },
            data: {
                Content: Content,
                URL: url,
                JSProgram: jsprogram,
            }
        };

        try {
            const response = await fetch(`https://webservice-381819.uc.r.appspot.com/package/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': '*/*',
                },
                body: JSON.stringify(requestData),
            });

            if (response.ok) {
                const responseData = await response.json();
                console.log(responseData);
            } else {
                const errorMessage = await response.text();
                console.error('Error:', errorMessage);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    };

    const handleClick = async (e) => {
        e.preventDefault();
        const isValid = await trigger(); // Trigger form validation for all fields

        if (isValid) {
            console.log('Errors:', errors);
            console.log('Calling handleSubmit');
            const formData = getValues();
            onSubmit(formData);
        } else {
            console.log('Form is not valid:', errors);
        }
    };




    return (
        <Container component="main" tabIndex="0">
            <Box mt={4} mb={4}>
                <Typography variant="h4" component="h1" gutterBottom>
                    Update Package
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
                                required
                                autoFocus
                                tabIndex="0"
                                {...register('id', { required: true })}
                                inputProps={{
                                    'aria-label': 'Package ID',
                                }}
                            />
                        </FormControl>
                    </Grid>
                    <Grid item xs={12}>
                        <FormControl fullWidth variant="outlined">
                            <InputLabel htmlFor="content">Content (ZIP File)</InputLabel>
                            <OutlinedInput
                                id="content"
                                type="file"
                                accept=".zip"
                                onChange={handleUpload}
                                label="Content (ZIP File)"
                                required
                                tabIndex="0"
                                inputProps={{
                                    ref: register('Content', { required: true }),
                                    'aria-label': 'Content (ZIP File)',
                                }}
                                startAdornment={
                                    <InputAdornment position="start">
                                        <IconButton edge="start" aria-label="Attach ZIP file">
                                            <AttachFile />
                                        </IconButton>
                                    </InputAdornment>
                                }
                            />
                        </FormControl>
                    </Grid>
                    <Grid item xs={12}>
                        <Button
                            variant="contained"
                            color="primary"
                            type="submit"
                            onClick={handleClick}
                            tabIndex="0"
                            aria-label="Submit package update"
                        >
                            Submit
                        </Button>
                        <input
                            type="hidden"
                            name="Name"
                            {...register('Name', { defaultValue: '' })}
                        />
                        <input
                            type="hidden"
                            name="Version"
                            {...register('Version', { defaultValue: '' })}
                        />
                        <input
                            type="hidden"
                            name="url"
                            {...register('url', { defaultValue: '' })}
                        />
                        <input
                            type="hidden"
                            name="jsprogram"
                            {...register('jsprogram', { defaultValue: '' })}
                        />

                    </Grid>
                </Grid>
            </form>
        </Container>
    );
};

export default UpdatePackage;
