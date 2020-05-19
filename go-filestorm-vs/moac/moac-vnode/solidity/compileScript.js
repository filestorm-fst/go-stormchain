

function compileSol(solFileName) {
  const solc = require('solc');
  const fs = require('fs');

  var abiRunFileName = './runFile.abi';
  var hexRunFileName = './runFile.bin';
  var wholeFileName = './runFile.all';

  console.log("Compiling....", solFileName);
  // 1) compile the solidity code
  var content = fs.readFileSync(solFileName, 'utf8');
  var output = solc.compile(content, 1);
  // console.log("output", output);


  for (var objContractName in output.contracts) {
    var contractJson = output.contracts[objContractName];
    var contractName = objContractName.substring(1);
    var abiFileName = './' + contractName + '.abi';
    var hexFileName = './' + contractName + '.bin';
    var fhFileName = './' + contractName + '.funcHash';
    // // console.dir(output.contracts);
    // fs.writeFileSync(wholeFileName, JSON.stringify(output));
    fs.writeFileSync(abiFileName, 'contractAbi =' + contractJson.interface + ';');
    fs.writeFileSync(hexFileName, 'contractBytecode = \'0x' + contractJson.bytecode + '\';');
    fs.writeFileSync(fhFileName, 'functionHashes =' + JSON.stringify(contractJson.functionHashes) + ' ;');
    console.log('contract name:', contractName);
    // console.log('contract abi:\n', contractJson.interface);
    // console.log('contract bytecode:\n', contractJson.bytecode);
    // console.log('contract function hashes:\n', JSON.stringify(contractJson.functionHashes));
  }
  // fs.writeFileSync(abiFileName, contractJson.interface);
  // fs.writeFileSync(hexFileName, contractJson.bytecode);
  // fs.writeFileSync(abiRunFileName, 'contractAbi =' + contractJson.interface + ';');
  // fs.writeFileSync(hexRunFileName, 'contractBytecode = \'' + contractJson.bytecode + '\';');

}

// console.log("Compiling....", process.argv[1]);
compileSol(process.argv[2]);