import App from './system'
import Page from './views/page'
import RPC from './system/rpc'
import Splash from './views/splash'

const wsURL = window.location.protocol.replace(/^http/, 'ws') + '//' + window.location.host + '/ws'
const app = new App(new RPC(wsURL), 'guest', '#content-wrapper')
const homePage = new Page('Вітаю!', new Splash(app))

app.setMenu('#navbarCollapse', `
		<div class="collapse navbar-collapse order-3" id="navbarCollapse">
			<ul class="navbar-nav">
				<li class="nav-item"><a href="/" data-navigo class="nav-link">Домівка</a></li>
			</ul>
		</div>`)
	.route('/', homePage)
