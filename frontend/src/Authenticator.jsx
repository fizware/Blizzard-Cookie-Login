import {CookieLogin} from "../wailsjs/go/main/App.js";
import {useState} from "react";
import {Alert, Box, Button, CircularProgress, Snackbar, TextField} from "@mui/material";
import "./assets/css/authenticator.css"

function Authenticator() {
    const [cookie, setCookie] = useState("");
    const [loading, setLoading] = useState(false);
    const [displayError, setDisplayError] = useState(false);
    const [error, setError] = useState("");

    const onCookieChange = (e) => {
        setCookie(e.target.value);
    }

    const onLoginClick = () => {
        if (loading) return;
        setLoading(true);
        CookieLogin(cookie).then(e => {
            if (e !== "") {
                setError(e);
                setDisplayError(true);
            }
            setLoading(false);
        });
    }

    const onErrorClose = () => {
        setDisplayError(false);
        setError("");
    }

    return (
        <div>
            <div className="center-div">
                <div>
                    <TextField className="cookie-box" variant="outlined" value={cookie} onChange={onCookieChange} label="Cookie" />
                </div>
                <Box className="login-div">
                    <Button disabled={loading} onClick={onLoginClick} className="cookie-login-button" variant="outlined">Login</Button>
                    {loading && (
                        <CircularProgress
                            size={24}
                            sx={{
                                position: 'absolute',
                                top: '50%',
                                left: '50%',
                                marginTop: '-12px',
                                marginLeft: '-12px',
                            }}
                        />
                    )}
                </Box>
            </div>
            <Snackbar anchorOrigin={{horizontal: "left", vertical: "bottom"}} open={displayError} autoHideDuration={1500} onClose={onErrorClose}>
                <Alert onClose={onErrorClose} severity="info" sx={{ width: '100%' }}>
                    {error}
                </Alert>
            </Snackbar>
        </div>
    )
}

export default Authenticator;