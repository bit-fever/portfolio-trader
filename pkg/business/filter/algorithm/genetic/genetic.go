//=============================================================================
/*
Copyright © 2024 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package genetic

import "github.com/bit-fever/portfolio-trader/pkg/business/filter/algorithm/optimization"

//=============================================================================

type geneticAlgorithm struct {
	ctx    optimization.Context
	config *optimization.GeneticConfig
}

//=============================================================================

func New() optimization.Algorithm {
	return &geneticAlgorithm{}
}

//=============================================================================
//===
//=== Genetic algorithm implementation
//===
//=============================================================================

func (ga *geneticAlgorithm) Init(ctx optimization.Context) {
	ga.ctx    = ctx
	ga.config = &ctx.AlgorithmConfig().Genetic

	//--- TODO Validate config and return an error in case of errors

	ga.config = &optimization.GeneticConfig{
		PopulationSize: 1000,
		MinSteps      : 1000,
	}
}

//=============================================================================

func (ga *geneticAlgorithm) StepsCount() uint {
	return	ga.config.MinSteps
}

//=============================================================================

func (ga *geneticAlgorithm) Optimize() {

}

//=============================================================================
