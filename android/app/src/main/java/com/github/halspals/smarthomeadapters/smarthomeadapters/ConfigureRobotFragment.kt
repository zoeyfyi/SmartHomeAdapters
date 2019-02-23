package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import kotlinx.android.synthetic.main.fragment_configure_robot.*

/**
 * WIP: The final (optional) screen of the robot registration wizard, where the user configures the robot parameters.
 * TODO FETCH CONFIGS FROM SERVER
 */
class ConfigureRobotFragment : Fragment() {

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_configure_robot, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        finish_button.setOnClickListener { _ -> setConfigParametersAndFinish() }

        // Set up the grid's adapter to display the configuration parameters requested
        parameter_grid.adapter = ParameterAdapter(
            view.context,
            (activity as RegisterRobotActivity).chosenUseCase.parameters)
    }

    /**
     * WIP: Sets the configuration parameters in the web server and finishes the registration wizard.
     * TODO read configuration input, make API call, finish only if appropriate
     */
    private fun setConfigParametersAndFinish() {
        activity?.finish()
    }
}
