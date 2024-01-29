import React from 'react'
import {createRoot} from 'react-dom/client'
import App from './App'
import { ThemeProvider, createTheme } from '@mui/material/styles';

const container = document.getElementById('root')

const root = createRoot(container)

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});

root.render(
    <React.StrictMode>
        <ThemeProvider theme={darkTheme}>
            <App/>
        </ThemeProvider>
    </React.StrictMode>
)
