import './ContentTarget.scss'
import {ReactElement} from 'react'

export default function ContentTarget(props: {
	children: ReactElement
	className?: string
}) {
	return (
		<div
			className={
				props.className === undefined
					? 'contentTargetDefault'
					: 'contentTargetDefault ' + props.className
			}
		>
			{props.children}
		</div>
	)
}
