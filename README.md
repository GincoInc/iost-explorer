# IOST Blockchain Explorer

### Prerequisites
* Golang 1.16.0 (or newer) is required to build this project
* Mongodb 4.0 (or newer) is required, We recommend Redis stable version
* node.js 8.10   (or newer) is required

## Installation
```bash
get repo first:
git clone git@github.com:GincoInc/iost-explorer.git
```

### Frontend
```bash
#  Change Directory
cd frontend/

# build
npm install

# run in test
npm run dev

# run production build
npm run build
```

### Backend
```bash
#  Change Directory
cd backend/

# build
make

# run
./explorer
```

### Mongo Cron Task
```bash
#  Change Directory
cd backend/task

# build
make

# run
./explorer-task
```

### Docker Compose (run all services)
```bash
cp ./frontend/config/index.docker-compose.js ./frontend/config/index.js
cp ./backend/config/config.docker-compose.json ./backend/config/config.json
make docker-build
docker-compose up
```

```bash
open http://localhost:8080
```

## Contribution

Contribution of any forms is appreciated!

Currently, our core tech team is working intensively to develop the first stable version and core block chain structure. We will accept pull request after the first stable release published.

If you have any questions, please email to team@iost.io

## Community & Resources

Make sure to check out these resources as well for more information and to keep up to date with all the latest news about IOST project and team.

[IOSToken on Reddit](https://www.reddit.com/r/IOStoken)

[Telegram](https://t.me/officialios)

[Twitter](https://twitter.com/IOStoken)

[Official website](https://iost.io)

## Disclaimer

- IOS Blockchain is unfinished and some parts are highly experimental. Use the code at your own risk.

- THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
