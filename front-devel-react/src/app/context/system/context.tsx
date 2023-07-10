import {ReactElement, createContext, useContext, useState} from 'react'

const defaultValues = {
	darkTheme: true,
	URIAPI: 'http://localhost:8080/',
	//Enable this handle if you need to change the theme
	//setDarkTheme: (value: boolean) => {},
}

//Initailize the context, this identify the type of the context for use de hook in the components
export const SystemContext = createContext(defaultValues)

//Create a hook to use the context
export const useSystemContext = () => {
	const context = useContext(SystemContext)
	if (!context)
		throw new Error(
			'useSystemContext must be used with a useSystemContextProvider binding the index.tsx file',
		)

	return context
}

//Create a provider for add the context to the main tree of the app
export function SystemContextProvider(props: {children: ReactElement}) {
	const [darkTheme, _] = useState(defaultValues.darkTheme)
	const URIAPI = defaultValues.URIAPI

	return (
		<SystemContext.Provider value={{darkTheme, URIAPI /*, setDarkTheme*/}}>
			{props.children}
		</SystemContext.Provider>
	)
}
