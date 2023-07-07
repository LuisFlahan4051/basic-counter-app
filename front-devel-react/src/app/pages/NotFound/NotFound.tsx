import './NotFound.scss'
import Img404 from './img/404 Error-rafiki color-fire.svg'

export default function NotFound() {
	return (
		<div className='NotFound'>
			<p className='NotFound_msg'>Not Found</p>

			<img className='NotFound_img' src={Img404} alt='Img404' />
		</div>
	)
}
