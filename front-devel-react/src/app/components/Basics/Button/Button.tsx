import {useState} from 'react'
import './Button.scss'

export default function Button(props: {
	title: string
	onClick?: any
	type?: 'button' | 'submit' | 'reset' | undefined
	theme: string // 'dark' | 'light'
	backgroundColor: string // css format: #000000 | colorName
	borderColor: string // css format: #000000 | colorName
	hoverBackgroundColor?: string // css format: #000000 | colorName
	hoverBorderColor?: string // css format: #000000 | colorName
	textColor?: string // css format: #000000 | colorName
}) {
	const [isHover, setIsHover] = useState(false)

	const style = {
		color: props.textColor === undefined ? '#000000' : `${props.textColor}`,
		borderColor: `${props.borderColor}`,
		backgroundColor: `${props.backgroundColor}`,
	}
	const styleHover = {
		color: props.textColor === undefined ? '#000000' : `${props.textColor}`,
		borderColor:
			props.hoverBackgroundColor === undefined
				? `${props.borderColor}`
				: `${props.hoverBorderColor}`,
		backgroundColor:
			props.hoverBackgroundColor === undefined
				? `${props.backgroundColor}`
				: `${props.hoverBackgroundColor}`,
	}

	return (
		<button
			className='basicButton'
			onClick={props.onClick === undefined ? () => {} : props.onClick}
			css-basic-button-theme={props.theme}
			style={isHover ? styleHover : style}
			onMouseEnter={() => setIsHover(true)}
			onMouseLeave={() => setIsHover(false)}
			type={props.type}
		>
			{props.title}
		</button>
	)
}
