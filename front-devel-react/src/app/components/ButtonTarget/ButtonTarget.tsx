import {Link} from 'react-router-dom'
import './ButtonTarget.scss'

export default function ButtonTarget(props: {
	title: string
	icon: string | null
	specificTheme: string
	to: string
	onClick: any
}) {
	return (
		<Link
			title={props.title}
			className='ButtonTarget'
			data-button-target-theme={props.specificTheme}
			onClick={props.onClick}
			to={props.to}
		>
			{props.icon !== null ? (
				<img src={props.icon} alt={props.title} className='ButtonTarget_icon' />
			) : (
				props.title
			)}
		</Link>
	)
}
