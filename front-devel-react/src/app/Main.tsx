import './Main.scss'
import Home from './pages/Home/Home'
import {BrowserRouter, Routes, Route, Navigate} from 'react-router-dom'
import NotFound from './pages/NotFound/NotFound'

function Main() {
	/*-------------------- Main Render ------------------------- */
	return (
		<div className='Main' data-global-theme={'dark'}>
			<BrowserRouter>
				<Routes>
					<Route path='/' element={<Navigate to='home' replace={true} />} />

					<Route path='home' element={<Home />} />

					<Route path='*' element={<NotFound />} />
				</Routes>
			</BrowserRouter>
		</div>
	)
}

export default Main
