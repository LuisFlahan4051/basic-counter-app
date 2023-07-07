import './ButtonCircle.scss'
import {useState} from 'react'
import {Link} from 'react-router-dom'

export default function ButtonCircle(props: {
	title: string
	icon: string | null
	color: string
	specificTheme: string
	to: string
	onClick: any
}) {
	const [isHover, setIsHover] = useState(false)

	const style = {
		borderColor: `${props.color}`,
	}
	const styleHover = {
		borderColor: `${props.color}`,
		backgroundColor: `${props.color}`,
	}

	return (
		<Link
			className='btn-circle'
			style={isHover ? styleHover : style}
			onMouseEnter={() => setIsHover(true)}
			onMouseLeave={() => setIsHover(false)}
			data-button-circle-theme={props.specificTheme}
			onClick={props.onClick}
			title={props.title}
			to={props.to}
		>
			{props.icon !== null ? (
				<img src={props.icon} alt={props.title} className='btn-circle_icon' />
			) : (
				props.title
			)}
		</Link>
	)
}
