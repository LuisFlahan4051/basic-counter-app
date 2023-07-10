import './Main.scss'
import Home from './pages/Home/Home'
import {BrowserRouter, Routes, Route, Navigate} from 'react-router-dom'
import NotFound from './pages/NotFound/NotFound'
import {useContext} from 'react'
import {SystemContext} from './context/system/context'

function Main() {
	/* -------------- CONTEXT --------------*/
	const {darkTheme} = useContext(SystemContext)

	/*-------------------- Main Render ------------------------- */
	return (
		<div className='Main' data-global-theme={darkTheme ? 'dark' : 'light'}>
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
