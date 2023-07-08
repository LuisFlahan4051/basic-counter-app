import './Home.scss'
import IconBackgroundDark from './imgs/LogoFondoDark.svg'
import IconBackground from './imgs/LogoFondoWite.svg'

function Home() {
	/* -------------- RENDER --------------*/

	// Just set the background
	function Background() {
		return (
			<div className='background'>
				<img
					src={true ? IconBackgroundDark : IconBackground}
					alt='icon background'
					className='background__icon'
				/>
			</div>
		)
	}

	function MainButtons() {
		return (
			<div className='mainButtons'>
				<button className='mainButtons__button'>Button</button>
			</div>
		)
	}

	return (
		//Main content in home
		<div className='display_home'>
			{/* this don't interfere with the main content */}
			<Background />

			{/* Main content */}
			<div className='Home'>
				<MainButtons />
			</div>
		</div>
	)
}

export default Home
