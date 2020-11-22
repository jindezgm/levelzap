<!--
 * @Author: jinde.zgm
 * @Date: 2020-11-22 14:43:05
 * @Descripttion: 
-->
## levelzaplevelzap 
implements leveled logging based on zap. The output log level can be modified through SetLevel().Log level from small to large, severity from low to high. levelzap uses zap log level by default, DEBUG(-1), INFO(0), WARN(1), ERROR(2)...   
You can use levelzap directly like this without any extra operations:   
```go    
    import log "github.com/jindezgm/levelzap"    
    ...    
    log.Debug("helloworld!", zap.String("level", "DEBUG"))    
    log.Info("helloworld!", zap.String("level", "INFO"))    
    log.Error("helloworld!", zap.String("level", "ERROR"))
```
The DEBUG log in the above example cannot be output because the default log level is INFO. Of course, can also write it like this:   
```go    
    import log "github.com/jindezgm/levelzap"    
    ...    
    log.V(log.DEBUG).Info("helloworld!", zap.String("level", "DEBUG"))    
    log.V(log.INFO).Info("helloworld!", zap.String("level", "INFO"))    
    log.V(log.ERROR).Info("helloworld!", zap.String("level", "ERROR"))
```
If the log level of zap does not meet the requirements, You can customize the log level like this:   
```go    
    import log "github.com/jindezgm/levelzap"    ...    
    log.SetLevel(4)    
    log.V(3).Info("helloworld!", zap.Int("level", 3))    
    log.V(4).Info("helloworld!", zap.Int("level", 4))    
    log.V(5).Info("helloworld!", zap.Int("level", 5))
```
Logs with level 3 cannot be output because the log level is set to 4.Levelzap has many arguments which can be set through flagset, For the default value, see loggingT.initDefault().levelzap provides the New interface to create different loggers for application which require multiple loggers.   
```go    
    import log "github.com/jindezgm/levelzap"    
    l1 := log.New()    
    l2 := log.New()    
    l1.V(log.DEBUG).Info("hellowworld!", zap.Int("logger", 1))    
    l2.V(log.INFO).Info("hellowworld!", zap.Int("logger", 2))
```