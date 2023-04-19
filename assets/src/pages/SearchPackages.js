import React, { useState } from 'react';
import {
    Container,
    Typography,
    Box,
    TextField,
    Button,
    CircularProgress,
    List,
    ListItem,
    ListItemText,
    FormControl,
    FormControlLabel,
    RadioGroup,
    Radio,
} from '@mui/material';

const SearchPackages = () => {
    const [search, setSearch] = useState('');
    const [results, setResults] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [searchType, setSearchType] = useState('name');

    const onSubmit = async (e) => {
        e.preventDefault();

        setIsLoading(true);
        try {
            let response;
            if (searchType === 'name') {
                response = await fetch(
                    `http://0.0.0.0:8080/package/byName/${encodeURIComponent(search)}`
                );
            } else {
                response = await fetch('http://0.0.0.0:8080/package/byRegEx', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        "RegEx": `"${search}"`
                    }),
                });
            }

            if (response.ok) {
                const data = await response.json();
                setResults(data);
            } else {
                console.error('Error:', response.statusText);
            }
        } catch (error) {
            console.error('Error:', error);
        }
        setIsLoading(false);
    };

    const renderResult = (result) => {
        if (searchType === 'name' && result.PackageMetadata) {
            return (
                <>
                    Name: {result.PackageMetadata.Name}, Version: {result.PackageMetadata.Version}
                    <br />
                    Action: {result.Action}, User: {result.User.name}
                </>
            );
        } else if (searchType === 'regex') {
            return (
                <>
                    Name: {result.Name}, Version: {result.Version}
                </>
            );
        }
    };

    return (
        <Container>
            <Box mt={4} mb={4}>
                <Typography variant="h4" component="h1" gutterBottom>
                    Search Packages
                </Typography>
            </Box>
            <form onSubmit={onSubmit}>
                <Box display="flex" alignItems="center" mb={4}>
                    <TextField
                        label="Search"
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                    />
                    <Button
                        variant="contained"
                        color="primary"
                        type="submit"
                        disabled={isLoading}
                        style={{ marginLeft: 8 }}
                    >
                        Search
                    </Button>
                </Box>
                <FormControl component="fieldset">
                    <RadioGroup
                        row
                        value={searchType}
                        onChange={(e) => {
                            setSearchType(e.target.value);
                            setResults([]); // Clear the search results
                        }}
                    >
                        <FormControlLabel
                            value="name"
                            control={<Radio />}
                            label="Search by Name"
                        />
                        <FormControlLabel
                            value="regex"
                            control={<Radio />}
                            label="Search by Regex"
                        />
                    </RadioGroup>
                </FormControl>
            </form>
            {isLoading ? (
                <CircularProgress />
            ) : (
                <List>
                    {results.map((result, index) => (
                        <ListItem key={index}>
                            <ListItemText
                                primary={renderResult(result)}
                            />
                        </ListItem>
                    ))}
                </List>
            )}
        </Container>
    );
};

export default SearchPackages;
