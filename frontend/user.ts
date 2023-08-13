import App from './system'
import Dashboard from './views/dashboard'
import Page from './views/page'
import RPC from './system/rpc'

const wsURL = window.location.protocol.replace(/^http/, 'ws') + '//' + window.location.host + '/ws'
const app = new App(new RPC(wsURL), 'user', '#content-wrapper')

const homePage = new Page('Мій кабінет', new Dashboard(app))

app.setMenu('#navbarCollapse',`
		<div class="collapse navbar-collapse order-3" id="navbarCollapse">
			<ul class="navbar-nav">
				<li class="nav-item"><a href="/" data-navigo class="nav-link">Домівка</a></li>
			</ul>
			<ul class="order-1 order-md-3 navbar-nav navbar-no-expand ml-auto">
				<li class="nav-item">
					<a class="nav-link" id="logout" href="/logout" role="button">Вихід</a>
				</li>
			</ul>
		</div>`)
	.route('/', homePage)
