# MATM-Wallet-Miner
MATM Wallet Miner is an open source project developed in GO.

__How does it work?__
The software generates hexadecimal sequences of 32 bytes each, which will create a key.
The key will then be controlled by the software using free nodes (https://rpc.ankr.com/eth) which will return the wallet balance.
If it is greater than 0 it will mark it as valid.

__How can I use it?__
At the moment the software is only available for Windows platforms.
To use it, simply go to the build folder and you will find the file already compiled.
Otherwise you can download the source code and compile it yourself but you will need a compiler.

__What are Threads and GoRutines?__
Threads (also called GoRoutine) are child processes of the main (Main routine).
The higher the amount you choose, the faster the program will generate new keys.

Too many Threads can create inconvenience to your computer,
if you are unsure about the right amount for you just type 0 when prompted and the program will automatically detect the best settings for you.

## ATTENTION
It is RECOMMENDED (NOT OBLIGATORY) to disable windows defender or other antivirus when the miner is running as they limit the usable resources of the computer.
IF YOU DO NOT TRUST THE SOFTWARE, IT CAN CONTINUE TO WORK EVEN WITH DEFENDER ACTIVATED.

## --> IF YOU CONTINUE TO DOUBT ABOUT RELIABILITY READ THE CODE YOURSELF <--
NOBODY forces you to use my miner.
