package odigosconditionalattributes

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

// NewFactory returns a new factory for the Resource processor.
func NewFactory() processor.Factory {
	return processor.NewFactory(
		component.MustNewType("odigosconditionalattributes"),
		createDefaultConfig,
		processor.WithTraces(createTracesProcessor, component.StabilityLevelBeta),
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createTracesProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Traces) (processor.Traces, error) {

	uniqueNewAttributes := calculateUniqueNewAttributes(cfg.(*Config))
	proc := &conditionalAttributesProcessor{logger: set.Logger, config: cfg.(*Config), uniqueNewAttributes: uniqueNewAttributes}

	return processorhelper.NewTraces(
		ctx,
		set,
		cfg,
		nextConsumer,
		proc.processTraces,
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}),
	)
}

func calculateUniqueNewAttributes(config *Config) map[string]struct{} {
	uniqueNewAttributes := make(map[string]struct{})
	for _, rule := range config.Rules {
		for _, value := range rule.NewAttributeValueConfigurations {
			for _, v := range value {
				if v.NewAttributeName != "" {
					uniqueNewAttributes[v.NewAttributeName] = struct{}{}
				}
			}
		}
	}
	return uniqueNewAttributes
}
