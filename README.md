# Telegram-ERC20-Factory-Bot

This is a relatively-simple telegram bot to collect data (name, symbol, total supply and type) about user's desired voting token, which then is passed to the website where user can approve the transaction using his Metamask wallet.


# Quick setup

1. Clone this repo
2. Create .secret file in the root directory, containing your individual Telegram bot API-key, acquired from Botfather
3. Execute npm install
4. Execute npm start
5. Execute go run .

Then you may communicate with your telegram bot which will provide you the link to the locally running website, where you may connect your Metamask wallet & mint your token!

# Note:

Factory smart-contract is currently deployed on Goerli testnet, so you'll need Goerli test eth to operate it.

Link to Factory contract on Etherscan: https://goerli.etherscan.io/address/0xAf5B8690245087a57128ec9543931574fDfAB4f1

Link to Factory contract's repo: https://github.com/daseinsucks/ERC20Factory


