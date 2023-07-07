import './AlertSceen.scss'
import iconOk from './icons/iconOk.svg'
import iconCancel from './icons/iconCancel.svg'
import iconQuestion from './icons/iconQuestion.svg'
import iconInfo from './icons/iconInfo.svg'
import iconWarning from './icons/iconWarning.svg'
import {useEffect, useRef} from 'react'

function AlertScreen(props: {
	type: string
	msg: string
	onAccept: () => void
	onCancel: () => void
	onClose: () => void
}) {
	const inputAccept = useRef(document.createElement('button'))
	const inputCancel = useRef(document.createElement('button'))

	let icon: any
	switch (props.type) {
		case 'Success':
			icon = iconOk
			break
		case 'Error':
			icon = iconCancel
			break
		case 'Info':
			icon = iconInfo
			break
		case 'Question':
			icon = iconQuestion
			break
		case 'Warning':
			icon = iconWarning
			break
		default:
			icon = iconOk
			break
	}

	useEffect(() => {
		inputAccept.current.focus()
	}, [])

	/* -------------- RENDER --------------*/
	return (
		<div className='AlertScreen'>
			<div className='alert__background' onClick={props.onClose}></div>

			<div className='alert__content'>
				<div className='alert__target'>
					<img src={icon} alt='alert__icon' />
					<p className='alert__msg'>{props.msg}</p>
					<div className='alert__buttons'>
						{props.type === 'Question' || props.type === 'Warning' ? (
							<>
								<button
									className='alert-btn alert__acept'
									onClick={props.onAccept}
									ref={inputAccept}
									onKeyDown={e => {
										if (e.key === 'ArrowRight') {
											inputCancel.current.focus()
										}
										if (e.key === 'Escape') {
											props.onCancel()
										}
									}}
								>
									Aceptar
								</button>
								<button
									className='alert-btn alert__cancel'
									onClick={props.onCancel}
									ref={inputCancel}
									onKeyDown={e => {
										if (e.key === 'ArrowLeft') {
											inputAccept.current.focus()
										}
										if (e.key === 'Escape') {
											props.onCancel()
										}
									}}
								>
									Cancelar
								</button>
							</>
						) : (
							<button
								className='alert-btn alert__ok'
								onClick={props.onAccept}
								ref={inputAccept}
								onKeyDown={e => {
									if (e.key === 'Escape') {
										props.onClose()
									}
								}}
							>
								OK
							</button>
						)}
					</div>
				</div>
			</div>
		</div>
	)
}

export default AlertScreen
