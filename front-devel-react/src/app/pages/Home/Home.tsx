import './Home.scss'
import IconBackgroundDark from './imgs/LogoFondoDark.svg'
import IconBackground from './imgs/LogoFondoWite.svg'
import Button from '../../components/Basics/Button/Button'
import {useContext, useState} from 'react'
import {SystemContext} from '../../context/system/context'
import ContentTarget from '../../components/Basics/ContentTarget/ContentTarget'
import {useForm} from 'react-hook-form'

function Home() {
	/* -------------- CONTEXT --------------*/
	const {darkTheme, URIAPI} = useContext(SystemContext)
	const [formState, setFormState] = useState(0)
	const {register, handleSubmit} = useForm()

	/* -------------- FUNCTIONS --------------*/

	function addIncomeAction() {
		console.log('addIncomeAction')
		setFormState(1)
	}

	async function incomeSubmit(data: any, event: any) {
		event.preventDefault()
		console.log('agregando')

		//for build the body of the request u need to use states.
		const response = await fetch(URIAPI + 'incomes', {
			method: 'GET',
			mode: 'cors',
			credentials: 'include',
			headers: {'Content-Type': 'application/json'},
			//body: JSON.stringify(data),
		})
		const {data1} = await response.json()

		console.log('Income added with id:' + data1.id)
	}

	function addExpenseAction() {
		console.log('addExpenseAction')
		setFormState(2)
	}

	/* -------------- RENDER --------------*/

	// Just set the background
	function Background() {
		return (
			<div className='background'>
				<img
					src={darkTheme ? IconBackgroundDark : IconBackground}
					alt='icon background'
					className='background__icon'
				/>
			</div>
		)
	}

	function RenderForm() {
		switch (formState) {
			case 1:
				return (
					<ContentTarget className='inputs__contentTarget'>
						<form
							onSubmit={handleSubmit(incomeSubmit)}
							className='inputs__form'
						>
							<h1>Nuevo Ingreso</h1>

							<label htmlFor='type'>Día:</label>
							<input
								type='date'
								id='day'
								defaultValue={new Date().toISOString().slice(0, 10)}
								{...register('created_at')}
							/>
							<label htmlFor='type'>Tipo:</label>
							<input type='text' id='type' {...register('type')} />
							<label htmlFor='cash'>Cantidad:</label>
							<input type='number' id='value' {...register('value')} />
							<Button
								theme={darkTheme ? 'dark' : 'light'}
								title='Agregar'
								backgroundColor='#407BFF'
								borderColor='#407BFF'
								hoverBackgroundColor='#6897fd'
								hoverBorderColor='#6897fd'
								type='submit'
							/>
							<Button
								theme={darkTheme ? 'dark' : 'light'}
								title='Limpiar'
								backgroundColor='#6BAEA7'
								borderColor='#6BAEA7'
								hoverBackgroundColor='#8fd1c9'
								hoverBorderColor='#8fd1c9'
								type='reset'
							/>
						</form>
					</ContentTarget>
				)
			case 2:
				return (
					<ContentTarget className='inputs__contentTarget'>
						<form onSubmit={() => {}} className='inputs__form'>
							<h1>Nuevo Egreso</h1>

							<label htmlFor='type'>Día:</label>
							<input
								type='date'
								name='created_at'
								id='day'
								defaultValue={new Date().toISOString().slice(0, 10)}
							/>
							<label htmlFor='type'>Tipo:</label>
							<input type='text' name='type' id='type' />
							<label htmlFor='cash'>Cantidad:</label>
							<input type='number' name='value' id='value' />
							<Button
								theme={darkTheme ? 'dark' : 'light'}
								title='Agregar'
								backgroundColor='#407BFF'
								borderColor='#407BFF'
								hoverBackgroundColor='#6897fd'
								hoverBorderColor='#6897fd'
								type='submit'
							/>
							<Button
								theme={darkTheme ? 'dark' : 'light'}
								title='Limpiar'
								backgroundColor='#6BAEA7'
								borderColor='#6BAEA7'
								hoverBackgroundColor='#8fd1c9'
								hoverBorderColor='#8fd1c9'
								type='reset'
							/>
						</form>
					</ContentTarget>
				)
			default:
				return <></>
		}
	}

	return (
		//Main content in home
		<div className='display_home'>
			{/* this don't interfere with the main content */}
			<Background />

			{/* Main content */}
			<div className='Home'>
				<div className='display_home__options'>
					<div className='options__content'>
						<Button
							title='Nuevo Ingreso'
							onClick={addIncomeAction}
							theme={darkTheme ? 'dark' : 'light'}
							backgroundColor='#178538'
							borderColor='#178538'
							hoverBackgroundColor='#1E9A4A'
							hoverBorderColor='#1E9A4A'
							textColor='#FFFFFF'
						/>
						<Button
							title='Nuevo Egreso'
							onClick={addExpenseAction}
							theme={darkTheme ? 'dark' : 'light'}
							backgroundColor='#ec3636'
							borderColor='#ec3636'
							hoverBackgroundColor='#f55656'
							hoverBorderColor='#f55656'
							textColor='#FFFFFF'
						/>
					</div>
				</div>

				<div className='display_home__inputs'>
					<RenderForm />
				</div>

				<div className='display_home__table'></div>
				<div className='display_home__graphic'></div>
			</div>
		</div>
	)
}

export default Home
