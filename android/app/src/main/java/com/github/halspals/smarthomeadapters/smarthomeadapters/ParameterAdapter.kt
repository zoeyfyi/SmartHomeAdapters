package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.support.design.card.MaterialCardView
import android.support.v7.widget.RecyclerView
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.SeekBar
import android.widget.Switch
import android.widget.TextView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.ConfigParameter
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.ConfigResult

/**
 * Provides an adapter to list [ConfigParameter]s in a [RecyclerView] of [MaterialCardView]s.
 *
 * @property parameters the configuration parameters to list
 */
class ParameterAdapter(
        private val parameters: List<ConfigParameter>
):  RecyclerView.Adapter<ParameterAdapter.ParameterAdapterViewHolder>() {

    override fun getItemCount(): Int = parameters.size

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ParameterAdapterViewHolder {
        // Otherwise inflate the appropriate card and set up its fields
        val cardView = LayoutInflater.from(parent.context)
                .inflate(R.layout.view_config_card, parent, false)
        return ParameterAdapterViewHolder(cardView as MaterialCardView)

    }

    override fun onBindViewHolder(viewHolder: ParameterAdapter.ParameterAdapterViewHolder, position: Int) {

        val parameter = parameters[position]
        val cardView = viewHolder.cardView
        val switch = cardView.findViewById<Switch>(R.id.input_switch)
        val seekBar = cardView.findViewById<SeekBar>(R.id.input_seekbar)
        val seekBarValueTextView = cardView.findViewById<TextView>(R.id.seekbar_value_text_view)

        when (parameter.type) {
            ConfigParameter.BOOL_TYPE -> {
                // The parameter is a boolean one -- set up the switch
                seekBar.visibility = View.GONE
                seekBarValueTextView.visibility = View.GONE
                switch.isChecked = parameter.details.current != 0
                switch.setOnCheckedChangeListener { _, b ->
                  parameter.details.current = if (b) { 1 } else { 0 }
                }
            }

            ConfigParameter.INT_TYPE -> {
                // The parameter is an integer one -- set up the seek bar
                switch.visibility = View.GONE

                val max = parameter.details.max
                val min = parameter.details.min
                val current = parameter.details.current

                seekBar.max = max - min // manually add min later as api <26 doesn't support seekBar.setMin()
                seekBar.progress = current - min
                seekBarValueTextView.text = "$current"

                // Set up the change listener for the seekbar
                seekBar.setOnSeekBarChangeListener(object : SeekBar.OnSeekBarChangeListener {
                    override fun onStopTrackingTouch(p0: SeekBar?) {}
                    override fun onStartTrackingTouch(p0: SeekBar?) {}

                    override fun onProgressChanged(seekBar: SeekBar?, progress: Int, byUser: Boolean) {
                        // add back the min we subtracted earlier as the seekbar starts at 0
                        val newValue = progress + min
                        seekBarValueTextView.text = newValue.toString()
                        parameter.details.current = newValue
                    }
                })
            }

            else -> {
                TODO("No other types expected")
            }
        }

        // Regardless of the type of card, set up the name and description fields
        cardView.findViewById<TextView>(R.id.config_name_text_view).text = parameter.name
        cardView.findViewById<TextView>(R.id.config_explanation_text_view).text = parameter.description
    }

    class ParameterAdapterViewHolder(internal val cardView: MaterialCardView) : RecyclerView.ViewHolder(cardView)

    internal fun getConfigResultsSnapshot(): List<ConfigResult> {
        return parameters.map { ConfigResult(it.id, it.type, it.details.current.toString()) }
    }

}