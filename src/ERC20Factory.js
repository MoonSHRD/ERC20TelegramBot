import React from "react";
import { useEffect, useState } from "react";
import {
  factoryContract,
  connectWallet,
  getCurrentWalletConnected,
  mintToken,
} from "./util/interact.js";

import logo from "./logo.svg";

const HelloWorld = () => {
  //state variables
  const [walletAddress, setWallet] = useState("");
  const [status, setStatus] = useState("");
  const [address, setAddress] = useState("Mint the token to begin."); //default message
  const [tokenName, setTokenName] = useState("");
  const [tokenTicker, setTokenTicker] = useState("");
  const [tokenDecimals, setTokenDecimals] = useState("");
  const [tokenType, setTokenType] = useState("");
  const queryParams = new URLSearchParams(window.location.search);



  //called only once
  useEffect(async () => {
    addSmartContractListener();
    const name = queryParams.get('name');
    const symbol = queryParams.get('symbol');
    const supply = queryParams.get('supply');
    const type = queryParams.get('type');

    setTokenName(name);
    setTokenTicker(symbol);
    setTokenDecimals(supply);
    setTokenType(type);
    
      async function fetchWallet() {
        const {address, status} = await getCurrentWalletConnected();
        setWallet(address);
        setStatus(status); 
      }
    fetchWallet();
    addWalletListener();
  
  }, []);


    function addSmartContractListener() {
      factoryContract.events.tokenCreated({}, (error, data) => {
        if (error) {
          setStatus("😥 " + error.message);
        } else {
          console.log(data.returnValues[0]);
          setAddress(data.returnValues[0]);
          setStatus("Token minted successfully!")
        }
      });
    }
    


    function addWalletListener() {
      if (window.ethereum) {
        window.ethereum.on("accountsChanged", (accounts) => {
          if (accounts.length > 0) {
            setWallet(accounts[0]);
            setStatus("You are ready to mint your token.");
          } else {
            setWallet("");
            setStatus("🦊 Connect to Metamask using the top right button.");
          }
        });
      } else {
        setStatus(
          <p>
            {" "}
            🦊{" "}
            <a target="_blank" href={`https://metamask.io/download.html`}>
              You must install Metamask, a virtual Ethereum wallet, in your
              browser.
            </a>
          </p>
        );
      }
    }

  const connectWalletPressed = async () => {
    const walletResponse = await connectWallet();
    setStatus(walletResponse.status);
    setWallet(walletResponse.address);
  };

  const onUpdatePressed = async () => {
    const { status } = await mintToken(walletAddress, tokenName, tokenTicker, tokenDecimals, tokenType);
    setStatus(status);
};

  //the UI of our component
  return (
    <div id="container">
      <img id="logo" src={logo}></img>
      <button id="walletButton" onClick={connectWalletPressed}>
        {walletAddress.length > 0 ? (
          "Connected: " +
          String(walletAddress).substring(0, 6) +
          "..." +
          String(walletAddress).substring(38)
        ) : (
          <span>Connect Wallet</span>
        )}
      </button>

      <h2 style={{ paddingTop: "50px" }}>Minted token address:</h2>
      <p>{address}</p>

      <h2 style={{ paddingTop: "18px" }}>Token details:</h2>

      <div>
        <input
          type="text"
          disabled = "true"
          placeholder="Token name."
          onChange={(e) => setTokenName(e.target.value)}
          value={tokenName}
        />

        <input
          type="text"
          disabled = "true"
          placeholder="Token ticker."
          onChange={(e) => setTokenTicker(e.target.value)}
          value={tokenTicker}
        />

        <input
          type="text"
          disabled = "true"
          placeholder="Token supply."
          onChange={(e) => setTokenDecimals(e.target.value)}
          value={tokenDecimals}
        />

        <input
          type="text"
          disabled = "true"
          placeholder="Token type."
          onChange={(e) => setTokenType(e.target.value)}
          value={tokenType}
        />


        <p id="status">{status}</p>

        <button id="publish" onClick={onUpdatePressed}>
          Mint token
        </button>
      </div>
    </div>
  );
};

export default HelloWorld;
