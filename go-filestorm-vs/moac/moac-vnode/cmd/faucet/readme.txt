for faucet to work:
1. start a moac with lightserv mode enabled. 
./moac --mine --lightserv 20

2. get key for the funding account

3. start faucet with following parameters.
./faucet -network 1 -genesis genesis.json -account.json key -account.pass "" -bootnodes "enode://8b792c0e6c4f875e7b8ddb6e687c9a54d498529108af5419b06a35fe504278d6e9a4216d5abba86f2ad85cd09d69db46c8cd7bb7b149ca650e361fc6bdc283ff@23.92.71.88:30333,enode://d28b17257bc683e3befd8d8b59f8e3b5e5829d451ad0ff0c7e5e2b760d516720ec24ecb75b4c5ff853dcbf5e82eeda3024ffef0de4e68900eb6869a52c4e2665@107.155.102.249:30333,enode://18bfa610fdcd52f9327c14e35ab2be8440eb1fcf530ff7f86cd3ffc9393f99485d1a55f72fc7e54b7b5aea12c3fda13b898ea613ca0ea2ffef435b07ba983578@107.155.106.71:30333,enode://65e5df6c89365031781aab97e8ea2ec05d48ff618d8b028f048e7bab0e4e0e0d6a088f61e2192b19917556d779b98ef65e494e080d762891c68a261e4de8b49c@127.0.0.1:30333"

